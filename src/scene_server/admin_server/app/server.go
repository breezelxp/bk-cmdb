/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"configcenter/src/auth/authcenter"
	"configcenter/src/common/backbone"
	"configcenter/src/common/backbone/configcenter"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/configures"
	svc "configcenter/src/scene_server/admin_server/service"
	"configcenter/src/scene_server/admin_server/synchronizer"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"github.com/emicklei/go-restful"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	process := new(MigrateServer)
	pconfig, err := configcenter.ParseConfigWithFile(op.ServConf.ExConfig)
	if nil != err {
		return fmt.Errorf("parse config file error %s", err.Error())
	}
	process.onHostConfigUpdate(*pconfig, *pconfig)
	service := svc.NewService(ctx)

	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     process.Config.Register.Address,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	process.Core = engine
	process.Service = service
	process.ConfigCenter = configures.NewConfCenter(ctx, engine.ServiceManageClient())
	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}
		db, err := local.NewMgo(process.Config.MongoDB.BuildURI(), 0)
		if err != nil {
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}
		process.Service.SetDB(db)
		process.Service.SetApiSrvAddr(process.Config.ProcSrvConfig.CCApiSrvAddr)
		err = process.ConfigCenter.Start(
			process.Config.Configures.Dir,
			process.Config.Errors.Res,
			process.Config.Language.Res,
		)
		if err != nil {
			return err
		}

		if process.Config.AuthCenter.Enable {
			authcli, err := authcenter.NewAuthCenter(nil, process.Config.AuthCenter)
			if err != nil {
				return fmt.Errorf("new authcenter client failed: %v", err)
			} else {
				process.Service.SetAuthcenter(authcli)
			}

			// synchronizer used for synchronizing resources between iam and cmdb
			sync := synchronizer.NewSynchronizer(ctx, &process.Config.AuthCenter, process.Core)
			blog.Info("begin to start synchronizer ...")
			err = sync.Run()
			if err != nil {
				return fmt.Errorf("run auth synchronizer routine failed %s", err.Error())
			}
		}
		break
	}
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(service.WebService())); err != nil {
		return err
	}
	<-ctx.Done()
	blog.V(0).Info("process stopped")
	return nil
}

type MigrateServer struct {
	Core         *backbone.Engine
	Config       *options.Config
	Service      *svc.Service
	ConfigCenter *configures.ConfCenter
}

var configLock sync.Mutex

func (h *MigrateServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigMap) > 0 {
		if h.Config == nil {
			h.Config = new(options.Config)
		}

		out, _ := json.MarshalIndent(current.ConfigMap, "", "  ") //ignore err, cause ConfigMap is map[string]string
		blog.V(3).Infof("config updated: \n%s", out)

		mongoConf := mongo.ParseConfigFromKV("mongodb", current.ConfigMap)
		h.Config.MongoDB = mongoConf

		h.Config.Errors.Res = current.ConfigMap["errors.res"]
		h.Config.Language.Res = current.ConfigMap["language.res"]
		h.Config.Configures.Dir = current.ConfigMap["confs.dir"]

		h.Config.Register.Address = current.ConfigMap["register-server.addrs"]

		h.Config.ProcSrvConfig.CCApiSrvAddr, _ = current.ConfigMap["procsrv.cc_api"]

		var err error
		h.Config.AuthCenter, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
		if err != nil && h.Config.AuthCenter.Enable {
			blog.Errorf("parse authcenter error: %v, config: %+v", err, current.ConfigMap)
		}
	}
}

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}
	return info, nil
}

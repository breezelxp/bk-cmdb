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

package service

import (
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	//	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) GetProcBindTemplate(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]

	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("params error :%v,appID:%v, rid:%s", err, appIDStr, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	procIDStr := pathParams[common.BKProcessIDField]
	procID, err := strconv.ParseInt(procIDStr, 10, 64)
	if nil != err {
		blog.Errorf("params error :%v,procID:%+v,rid:%s", err, procIDStr, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	// search object instance
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	input := new(meta.QueryCondition)
	input.Condition = condition

	tempRet, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDConfigTemp, input)
	if err != nil || !tempRet.Result {
		blog.Errorf("fail to GetProcBindTemplate when do searchobject. err:%v, errcode:%d, errmsg:%s,rid:%s", err, tempRet.Code, tempRet.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrObjectSelectInstFailed)})
		return
	}

	condition[common.BKProcessIDField] = procID

	// get process to templation by condition
	proc2TempRet, err := ps.CoreAPI.ProcController().SearchProc2Template(srvData.ctx, srvData.header, condition)
	if err != nil || !proc2TempRet.Result {
		blog.Errorf("fail to GetProcessTemplate when do GetProc2Template. err:%v, errcode:%d, errmsg:%s,rid:%s", err, proc2TempRet.Code, proc2TempRet.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcSelectBindToMoudleFaile)})
		return
	}

	result := make([]interface{}, 0)

	for _, temp := range tempRet.Data.Info {
		iTempID, err := util.GetInt64ByInterface(temp[common.BKTemlateIDField])
		if nil != err {
			continue
		}
		sTempName, ok := temp[common.BKTemplateNameField].(string)
		if false == ok {
			continue
		}
		sFileName, ok := temp[common.BKFileNameField].(string)
		if false == ok {
			continue
		}
		isBind := 0
		for _, proc2Temp := range proc2TempRet.Data {
			jTempID, err := util.GetInt64ByInterface(proc2Temp[common.BKTemlateIDField])
			if nil != err {
				continue
			}
			if iTempID == jTempID {
				isBind = 1
			}
		}
		result = append(result, types.MapStr{common.BKTemplateNameField: sTempName, common.BKFileNameField: sFileName, "is_bind": isBind})
	}
	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) BindProc2Template(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	pathParams := req.PathParameters()

	appIDStr := pathParams[common.BKAppIDField]
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("params error :%v,input:%+v,rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	procID, err := strconv.ParseInt(pathParams[common.BKProcessIDField], 10, 64)
	if nil != err {
		blog.Errorf("params error. err:%v, input:%#v, rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	tempID, err := strconv.ParseInt(pathParams[common.BKTemlateIDField], 10, 64)
	if nil != err {
		blog.Errorf("params error. err:%v, input:%+v, rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	params := make([]interface{}, 0)
	cell := make(map[string]interface{})
	cell[common.BKAppIDField] = appID
	cell[common.BKProcessIDField] = procID
	cell[common.BKTemlateIDField] = tempID
	cell[common.BKOwnerIDField] = util.GetOwnerID(req.Request.Header)
	params = append(params, cell)

	ret, err := ps.CoreAPI.ProcController().CreateProc2Template(srvData.ctx, srvData.header, params)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("fail to BindModuleProcess. err: %v, errcode:%d, errmsg: %s,input:%+v,condtion:%+v,rid:%s", err, ret.Code, ret.ErrMsg, pathParams, params, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcBindToTemplateFailed)})
		return
	}

	// save operation log
	log := meta.SaveAuditLogParams{
		ID:      int64(procID),
		Model:   common.BKInnerObjIDProc,
		Content: meta.Content{},
		OpDesc:  fmt.Sprintf("bind template [%d] to process [%d]", tempID, procID),
		OpType:  auditoplog.AuditOpTypeAdd,
		BizID:   int64(appID),
	}
	auditrsp, err := ps.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, log)
	if err != nil || (auditrsp != nil && !auditrsp.Result) {
		blog.Errorf("save auditlog failed %v %v,rid:%s", auditrsp, err, srvData.rid)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteProc2Template(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	pathParams := req.PathParameters()
	appIDStr := pathParams[common.BKAppIDField]
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("params error :%v,input:%+v,rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	procID, err := strconv.ParseInt(pathParams[common.BKProcessIDField], 10, 64)
	if nil != err {
		blog.Errorf("params error :%v,input:%+v,rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	tempID, err := strconv.ParseInt(pathParams[common.BKTemlateIDField], 10, 64)
	if nil != err {
		blog.Errorf("params error :%v, input:%+v,rid:%s", err, pathParams, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	condition[common.BKProcessIDField] = procID
	condition[common.BKTemlateIDField] = tempID

	ret, err := ps.CoreAPI.ProcController().DeleteProc2Template(srvData.ctx, srvData.header, condition)
	if err != nil {
		blog.Errorf("fail to delete process template bind. err: %v,  input:%#v, condition:%#v, rid:%s", err, pathParams, condition, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("fail to delete process template bind.  errcode:%s, errmsg: %s, input:%#v, condition:%#v, rid:%s", ret.Code, ret.ErrMsg, pathParams, condition, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	// save operation log
	log := meta.SaveAuditLogParams{
		ID:      int64(procID),
		Model:   common.BKInnerObjIDProc,
		Content: meta.Content{},
		OpDesc:  fmt.Sprintf("unbind template [%d] to process [%d]", tempID, procID),
		OpType:  auditoplog.AuditOpTypeDel,
		BizID:   int64(appID),
	}
	auditrsp, err := ps.CoreAPI.CoreService().Audit().SaveAuditLog(srvData.ctx, srvData.header, log)
	if err != nil || (auditrsp != nil && !auditrsp.Result) {
		blog.Errorf("save auditlog failed %v %v,rid:%s", auditrsp, err, srvData.rid)
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

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

package instance

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

type InstanceClientInterface interface {
	CreateInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelInstance) (resp *metadata.CreatedOneOptionResult, err error)
	CreateManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateManyModelInstance) (resp *metadata.CreatedManyOptionResult, err error)
	SetManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.SetManyModelInstance) (resp *metadata.SetOptionResult, err error)
	UpdateInstance(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	ReadInstance(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (resp *metadata.QueryConditionResult, err error)
	DeleteInstance(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	DeleteInstanceCascade(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	//  ReadInstanceStruct 按照结构体返回实例数据
	ReadInstanceStruct(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition,
		result interface{}) (err errors.CCErrorCoder)

	// CountInstances counts model instances num.
	CountInstances(ctx context.Context, header http.Header,
		objID string, input *metadata.Condition) (*metadata.CountResponse, error)
}

func NewInstanceClientInterface(client rest.ClientInterface) InstanceClientInterface {
	return &instance{client: client}
}

type instance struct {
	client rest.ClientInterface
}

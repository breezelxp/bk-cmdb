/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package process

import (
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateServiceInstance(ctx core.ContextParams, instance metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	// base attribute validate
	if field, err := instance.Validate(); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, instance.Metadata); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// keep metadata clean
	instance.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))

	// validate service template id field
	serviceTemplate, err := p.GetServiceTemplate(ctx, instance.ServiceTemplateID)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, service_template_id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "service_category_id")
	}

	// validate module id field
	if err = p.validateModuleID(ctx, instance.ModuleID); err != nil {
		blog.Errorf("CreateServiceInstance failed, module id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "module_id")
	}

	// validate host id field
	if err = p.validateHostID(ctx, instance.HostID); err != nil {
		blog.Errorf("CreateServiceInstance failed, host id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "host_id")
	}

	// make sure biz id identical with service template
	serviceTemplateBizID, err := metadata.BizIDFromMetadata(serviceTemplate.Metadata)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, parse biz id from service template failed, code: %d, err: %+v, rid: %s", common.CCErrCommInternalServerError, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParseBizIDFromMetadataInDBFailed)
	}
	if bizID != serviceTemplateBizID {
		blog.Errorf("CreateServiceInstance failed, validation failed, input bizID:%d not equal service template bizID:%d, rid: %s", bizID, serviceTemplateBizID, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameProcessTemplate)
	if nil != err {
		blog.Errorf("CreateServiceInstance failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, err
	}
	instance.ID = int64(id)

	instance.Creator = ctx.User
	instance.Modifier = ctx.User
	instance.CreateTime = time.Now()
	instance.LastTime = time.Now()
	instance.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Insert(ctx.Context, &instance); nil != err {
		blog.Errorf("CreateServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return nil, err
	}
	return &instance, nil
}

func (p *processOperation) GetServiceInstance(ctx core.ContextParams, templateID int64) (*metadata.ServiceInstance, error) {
	instance := metadata.ServiceInstance{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).One(ctx.Context, &instance); nil != err {
		blog.Errorf("GetServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		if err.Error() == "document not found" {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, err
	}

	return &instance, nil
}

func (p *processOperation) UpdateServiceInstance(ctx core.ContextParams, instanceID int64, input metadata.ServiceInstance) (*metadata.ServiceInstance, error) {
	instance, err := p.GetServiceInstance(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	if field, err := input.Validate(); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// update fields to local object
	// TODO: fixme with update other fields than name
	instance.Name = input.Name

	// do update
	filter := map[string]int64{common.BKFieldID: instanceID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Update(ctx, filter, instance); nil != err {
		blog.Errorf("UpdateServiceTemplate failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return nil, err
	}
	return instance, nil
}

func (p *processOperation) ListServiceInstance(ctx core.ContextParams, bizID int64, serviceTemplateID int64, hostID int64, limit metadata.BasePage) (*metadata.MultipleServiceInstance, error) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
	filter := map[string]interface{}{}
	filter["metadata"] = md.ToMapStr()

	if serviceTemplateID != 0 {
		filter["service_template_id"] = serviceTemplateID
	}

	if hostID != 0 {
		filter["host_id"] = hostID
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return nil, err
	}
	instances := make([]metadata.ServiceInstance, 0)
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Start(
		uint64(limit.Start)).Limit(uint64(limit.Limit)).All(ctx.Context, &instances); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return nil, err
	}

	result := &metadata.MultipleServiceInstance{
		Count: total,
		Info:  instances,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceInstance(ctx core.ContextParams, serviceInstanceID int64) error {
	instance, err := p.GetServiceInstance(ctx, serviceInstanceID)
	if err != nil {
		blog.Errorf("DeleteServiceInstance failed, GetServiceInstance failed, instanceID: %d, err: %+v, rid: %s", serviceInstanceID, err, ctx.ReqID)
		return err
	}

	// service template that referenced by process template shouldn't be removed
	usageFilter := map[string]int64{"service_instance_id": instance.ID}
	usageCount, err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(usageFilter).Count(ctx.Context)
	if nil != err {
		blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, ctx.ReqID)
		return err
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceInstance failed, forbidden delete service instance be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	deleteFilter := map[string]int64{common.BKFieldID: instance.ID}
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceTemplate, err, ctx.ReqID)
		return err
	}
	return nil
}
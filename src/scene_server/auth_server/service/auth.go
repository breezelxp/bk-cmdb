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
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

// Authorize works to check if a user has the authority to operate resources
func (s *AuthService) Authorize(ctx *rest.Contexts) {
	if !auth.EnableAuthorize() {
		ctx.RespEntity(types.Decision{Authorized: true})
		return
	}

	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	blog.InfoJSON("\n-> authorize request: %s, rid: %s\n", authAttribute, ctx.Kit.Rid)

	one := authAttribute.Resources[0]
	if one.Action == meta.SkipAction {
		ctx.RespEntity(meta.Decision{Authorized: true})
		return
	}

	// TODO: remove this after debug.
	if skip(one.Type) {
		ctx.RespEntity(meta.Decision{Authorized: true})
		return
	}

	action, resources, err := iam.AdaptAuthOptions(&one)
	if err != nil {
		blog.Errorf("Adaptor failed, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if action == iam.Skip {
		ctx.RespEntity(meta.Decision{Authorized: true})
		return
	}

	ops := &types.AuthOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   authAttribute.User.UserName,
		},
		Action: types.Action{
			ID: string(action),
		},
		Resources: []types.Resource{*resources},
	}

	decision, err := s.authorizer.Authorize(ctx.Kit.Ctx, ops)
	if err != nil {
		blog.ErrorJSON("Authorize failed, err: %s, ops: %s, authAttribute: %s, rid: %s", err, ops, authAttribute, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decision)
}

// AuthorizeBath works to check if a user has the authority to operate resources.
func (s *AuthService) AuthorizeBatch(ctx *rest.Contexts) {
	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// TODO: remove this after debug
	ctx.RespEntity(meta.Decision{Authorized: true})
	return

	decisions, err := s.authorizeBatch(ctx.Kit, authAttribute, true)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decisions)
}

// AuthorizeAnyBatch works to check if a user has any authority for actions.
func (s *AuthService) AuthorizeAnyBatch(ctx *rest.Contexts) {
	authAttribute := new(meta.AuthAttribute)
	err := ctx.DecodeInto(authAttribute)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	blog.InfoJSON("-> authorize any request: %s, rid: %s", authAttribute, ctx.Kit.Rid)

	decisions, err := s.authorizeBatch(ctx.Kit, authAttribute, false)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(decisions)
}

func skip(typ meta.ResourceType) bool {
	if typ != meta.Business &&
		typ != meta.DynamicGrouping &&
		typ != meta.ProcessServiceCategory &&
		typ != meta.EventPushing {
		return true
	}
	return false
}

func (s *AuthService) authorizeBatch(kit *rest.Kit, attr *meta.AuthAttribute, exact bool) ([]types.Decision, error) {
	if !auth.EnableAuthorize() {
		decisions := make([]types.Decision, len(attr.Resources))
		for i := range decisions {
			decisions[i].Authorized = true
		}
		return decisions, nil
	}

	authBatchArr := make([]*types.AuthBatch, 0)
	decisions := make([]types.Decision, len(attr.Resources))
	for index, resource := range attr.Resources {

		// TODO: remove this after debug.
		if skip(resource.Type) {
			decisions[index].Authorized = true
			continue
		}

		if resource.Action == meta.SkipAction {
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, kit.Rid)
			continue
		}

		action, resources, err := iam.AdaptAuthOptions(&resource)
		if err != nil {
			blog.Errorf("adaptor cmdb resource to iam failed, err: %s, rid: %s", err, kit.Rid)
			return nil, err
		}

		if action == iam.Skip {
			// this resource should be skipped, do not need to verify in auth center.
			decisions[index].Authorized = true
			blog.V(5).Infof("skip authorization for resource: %+v, rid: %s", resource, kit.Rid)
			continue
		}

		authBatchArr = append(authBatchArr, &types.AuthBatch{
			Action:    types.Action{ID: string(action)},
			Resources: []types.Resource{*resources},
		})
	}

	// all resources are skipped
	if len(authBatchArr) == 0 {
		return decisions, nil
	}

	ops := &types.AuthBatchOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   attr.User.UserName,
		},
		Batch: authBatchArr,
	}
	var authDecisions []*types.Decision
	var err error
	if exact {
		authDecisions, err = s.authorizer.AuthorizeBatch(kit.Ctx, ops)
	} else {
		authDecisions, err = s.authorizer.AuthorizeAnyBatch(kit.Ctx, ops)
	}
	if err != nil {
		blog.ErrorJSON("authorize batch failed, err: %s, ops: %s, auth: %s, rid: %s", err, ops, attr, kit.Rid)
		return nil, err
	}
	index := 0
	for _, decision := range authDecisions {
		// skip resources' decisions are already set as authorized
		for decisions[index].Authorized {
			index++
		}
		decisions[index].Authorized = decision.Authorized
		index++
	}
	return decisions, nil
}

// ListAuthorizedResources returns all specified resources the user has the authority to operate.
func (s *AuthService) ListAuthorizedResources(ctx *rest.Contexts) {
	input := new(meta.ListAuthorizedResourcesParam)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	iamResourceType, err := iam.ConvertResourceType(input.ResourceType, 0)
	if err != nil {
		blog.Errorf("ConvertResourceType failed, err: %+v, resourceType: %s, rid: %s", err, input.ResourceType, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	iamActionID, err := iam.ConvertResourceAction(input.ResourceType, input.Action, input.BizID)
	if err != nil {
		blog.ErrorJSON("ConvertResourceAction failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	resources := make([]types.Resource, 0)
	if input.BizID > 0 {
		businessPath := "/" + string(iam.Business) + "," + strconv.FormatInt(input.BizID, 10) + "/"
		resource := types.Resource{
			System: iam.SystemIDCMDB,
			Type:   types.ResourceType(*iamResourceType),
			Attribute: map[string]interface{}{
				types.IamPathKey: []string{businessPath},
			},
		}
		resources = append(resources, resource)
	}

	ops := &types.AuthOptions{
		System: iam.SystemIDCMDB,
		Subject: types.Subject{
			Type: "user",
			ID:   input.UserName,
		},
		Action: types.Action{
			ID: string(iamActionID),
		},
		Resources: resources,
	}
	resourceIDs, err := s.authorizer.ListAuthorizedInstances(ctx.Kit.Ctx, ops, types.ResourceType(*iamResourceType))
	if err != nil {
		blog.ErrorJSON("ListAuthorizedInstances failed, err: %+v, input: %s, ops: %s, input: %s, rid: %s", err, input, ops, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resourceIDs)
}

// GetNoAuthSkipUrl returns the redirect url to iam for user to apply for specific authorizations
func (s *AuthService) GetNoAuthSkipUrl(ctx *rest.Contexts) {
	input := new(metadata.IamPermission)
	err := ctx.DecodeInto(input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	url, err := esb.EsbClient().IamSrv().GetNoAuthSkipUrl(ctx.Kit.Ctx, ctx.Kit.Header, *input)
	if err != nil {
		blog.ErrorJSON("GetNoAuthSkipUrl failed, err: %s, input: %s, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(url)
}

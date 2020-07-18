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

package iam

import (
	"fmt"
	"strconv"

	"configcenter/src/ac/meta"
)

func GenerateResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	switch attribute.Basic.Type {
	case meta.Business:
		return businessResourceID(resourceType, attribute)
	case meta.Model:
		return modelResourceID(resourceType, attribute)
	case meta.ModelModule:
		return modelModuleResourceID(resourceType, attribute)
	case meta.ModelSet:
		return modelSetResourceID(resourceType, attribute)
	case meta.MainlineModel:
		return mainlineModelResourceID(resourceType, attribute)
	case meta.MainlineModelTopology:
		return mainlineModelTopologyResourceID(resourceType, attribute)
	case meta.MainlineInstanceTopology:
		return mainlineInstanceTopologyResourceID(resourceType, attribute)
	case meta.AssociationType:
		return associationTypeResourceID(resourceType, attribute)
	case meta.ModelAssociation:
		return modelAssociationResourceID(resourceType, attribute)
	case meta.ModelInstanceAssociation:
		return modelInstanceAssociationResourceID(resourceType, attribute)
	case meta.ModelInstance, meta.MainlineInstance:
		return modelInstanceResourceID(resourceType, attribute)
	case meta.ModelInstanceTopology:
		return modelInstanceTopologyResourceID(resourceType, attribute)
	case meta.ModelTopology:
		return modelTopologyResourceID(resourceType, attribute)
	case meta.ModelClassification:
		return modelClassificationResourceID(resourceType, attribute)
	case meta.ModelAttributeGroup:
		return modelAttributeGroupResourceID(resourceType, attribute)
	case meta.ModelAttribute:
		return modelAttributeResourceID(resourceType, attribute)
	case meta.ModelUnique:
		return modelUniqueResourceID(resourceType, attribute)
	case meta.UserCustom:
		return hostUserCustomResourceID(resourceType, attribute)
	case meta.HostFavorite:
		return hostFavoriteResourceID(resourceType, attribute)
	case meta.NetDataCollector:
		return netDataCollectorResourceID(resourceType, attribute)
	case meta.EventPushing:
		return eventSubscribeResourceID(resourceType, attribute)
	case meta.HostInstance:
		return hostInstanceResourceID(resourceType, attribute)
	case meta.DynamicGrouping:
		return dynamicGroupingResourceID(resourceType, attribute)
	case meta.AuditLog:
		return auditLogResourceID(resourceType, attribute)
	case meta.SystemBase:
		return make([]RscTypeAndID, 0), nil
	case meta.Plat:
		return platID(resourceType, attribute)
	case meta.Process:
		return processResourceID(resourceType, attribute)
	case meta.ProcessServiceInstance:
		return processServiceInstanceResourceID(resourceType, attribute)
	case meta.BizTopology:
		return bizTopologyResourceID(resourceType, attribute)
	case meta.ProcessTemplate:
		return processTemplateResourceID(resourceType, attribute)
	case meta.ProcessServiceCategory:
		return processServiceCategoryResourceID(resourceType, attribute)
	case meta.ProcessServiceTemplate:
		return processServiceTemplateResourceID(resourceType, attribute)
	case meta.SetTemplate:
		return setTemplateResourceID(resourceType, attribute)
	case meta.OperationStatistic:
		return operationStatisticResourceID(resourceType, attribute)
	case meta.HostApply:
		return hostApplyResourceID(resourceType, attribute)
	case meta.ResourcePoolDirectory:
		return resourcePoolDirectoryResourceID(resourceType, attribute)
	case meta.EventWatch:
		return resourceWatch(resourceType, attribute)
	case meta.CloudAccount:
		return cloudAccountResourceID(resourceType, attribute)
	case meta.CloudResourceTask:
		return cloudResourceTaskResourceID(resourceType, attribute)
	case meta.ConfigAdmin:
		return configAdminResourceID(resourceType, attribute)
	}
	return nil, fmt.Errorf("gen id failed: unsupported resource type: %s", attribute.Type)
}

// generate business related resource id.
func businessResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return make([]RscTypeAndID, 0), nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}

	return []RscTypeAndID{id}, nil
}

// generate model's resource id, works for app model and model management
// resource type in auth center.
func modelResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return make([]RscTypeAndID, 0), nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
	}
	id.ResourceID = strconv.FormatInt(attribute.InstanceID, 10)

	return []RscTypeAndID{id}, nil
}

// generate module resource id.
func modelModuleResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func modelSetResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func mainlineModelResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func mainlineModelTopologyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func mainlineInstanceTopologyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func modelAssociationResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func associationTypeResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return make([]RscTypeAndID, 0), nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}

	return []RscTypeAndID{id}, nil
}

func modelInstanceAssociationResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func modelInstanceResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		if len(attribute.Layers) == 0 {
			return make([]RscTypeAndID, 0), nil
		}
		// for create
		return []RscTypeAndID{{
			ResourceType: SysInstanceModel,
			ResourceID:   attribute.Layers[0].InstanceIDEx,
		}}, nil
	}

	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}

	return []RscTypeAndID{
		{
			ResourceType: SysInstanceModel,
			ResourceID:   attribute.Layers[0].InstanceIDEx,
		},
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func modelInstanceTopologyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func modelTopologyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func modelClassificationResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return make([]RscTypeAndID, 0), nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}
	return []RscTypeAndID{id}, nil
}

func modelAttributeGroupResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func modelAttributeResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	if attribute.BusinessID > 0 {
		id.ResourceType = BizCustomField
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func modelUniqueResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func hostUserCustomResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func hostFavoriteResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func netDataCollectorResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func hostInstanceResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	// translate all parent layers
	resourceIDs := make([]RscTypeAndID, 0)

	if attribute.InstanceID == 0 {
		return resourceIDs, nil
	}

	for _, layer := range attribute.Layers {
		iamResourceType, err := ConvertResourceType(layer.Type, attribute.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("convert resource type to iam resource type failed, layer: %+v, err: %+v", layer, err)
		}
		resourceID := RscTypeAndID{
			ResourceType: *iamResourceType,
			ResourceID:   strconv.FormatInt(layer.InstanceID, 10),
		}
		resourceIDs = append(resourceIDs, resourceID)
	}

	// append host resource id to end
	hostResourceID := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}
	resourceIDs = append(resourceIDs, hostResourceID)

	return resourceIDs, nil
}

func eventSubscribeResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func dynamicGroupingResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 && len(attribute.InstanceIDEx) == 0 {
		return make([]RscTypeAndID, 0), nil
	}

	instanceID := strconv.FormatInt(attribute.InstanceID, 10)
	if len(attribute.InstanceIDEx) != 0 {
		instanceID = attribute.InstanceIDEx
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   instanceID,
		},
	}, nil
}

func auditLogResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.InstanceIDEx) == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	instanceID := attribute.InstanceIDEx
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   instanceID,
		},
	}, nil
}

func platID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}

	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func processResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func processServiceInstanceResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func bizTopologyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func processTemplateResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func processServiceCategoryResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

func processServiceTemplateResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func setTemplateResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func operationStatisticResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func hostApplyResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func resourcePoolDirectoryResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func resourceWatch(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return make([]RscTypeAndID, 0), nil
}

func cloudAccountResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func cloudResourceTaskResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID == 0 {
		return make([]RscTypeAndID, 0), nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func configAdminResourceID(resourceType TypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return make([]RscTypeAndID, 0), nil
}

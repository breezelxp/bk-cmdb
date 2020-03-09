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

package validator

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

func getString(val interface{}) string {
	if val == nil {
		return ""
	}
	if ret, ok := val.(string); ok {
		return ret
	}
	return ""
}
func getBool(val interface{}) bool {
	if val == nil {
		return false
	}
	if ret, ok := val.(bool); ok {
		return ret
	}
	return false
}

// FillLostedFieldValue fill the value in inst map data
func FillLostedFieldValue(valData map[string]interface{}, propertys []metadata.Attribute, ignorefields []string) {
	ignores := map[string]bool{}
	for _, field := range ignorefields {
		ignores[field] = true
	}
	for _, field := range propertys {
		if field.PropertyID == common.BKChildStr || field.PropertyID == common.BKParentStr {
			continue
		}
		if ignores[field.PropertyID] {
			continue
		}
		_, ok := valData[field.PropertyID]
		if !ok {
			switch field.PropertyType {
			case common.FieldTypeSingleChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeLongChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeInt:
				valData[field.PropertyID] = nil
			case common.FieldTypeEnum:
				// parse enum error. not set default value
				enumOptions, _ := ParseEnumOption(field.Option)
				if len(enumOptions) > 0 {
					var defaultOption *EnumVal
					for _, k := range enumOptions {
						if k.IsDefault {
							defaultOption = &k
							break
						}
					}
					if nil != defaultOption {
						valData[field.PropertyID] = defaultOption.ID
					} else {
						valData[field.PropertyID] = nil
					}
				} else {
					valData[field.PropertyID] = nil
				}
			case common.FieldTypeDate:
				valData[field.PropertyID] = nil
			case common.FieldTypeTime:
				valData[field.PropertyID] = nil
			case common.FieldTypeUser:
				valData[field.PropertyID] = nil
			case common.FieldTypeMultiAsst:
				valData[field.PropertyID] = nil
			case common.FieldTypeTimeZone:
				valData[field.PropertyID] = nil
			case common.FieldTypeBool:
				valData[field.PropertyID] = false
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(val interface{}) (EnumOption, error) {
	enumOptions := []EnumVal{}
	if nil == val || "" == val {
		return enumOptions, nil
	}
	switch options := val.(type) {
	case []EnumVal:
		return options, nil
	case string:
		err := json.Unmarshal([]byte(options), &enumOptions)
		if nil != err {
			blog.Errorf("ParseEnumOption error : %s", err.Error())
			return nil, errors.New("format error. " + err.Error())
		}
	default:
		byteArr, err := json.Marshal(options)
		if nil != err {
			blog.Errorf("ParseEnumOption Marshal error : %s", err.Error())
			return nil, errors.New("marshal options error. " + err.Error())
		}
		err = json.Unmarshal(byteArr, &enumOptions)
		if nil != err {
			blog.Errorf("ParseEnumOption Unmarshal error : %s", err.Error())
			return nil, errors.New("unmarshal options error. " + err.Error())
		}

	}
	return enumOptions, nil
}

//parseMinMaxOption  parse int data in option
func parseMinMaxOption(val interface{}) MinMaxOption {
	minMaxOption := MinMaxOption{}
	if nil == val || "" == val {
		return minMaxOption
	}
	switch option := val.(type) {
	case string:

		minMaxOption.Min = gjson.Get(option, "min").Raw
		minMaxOption.Max = gjson.Get(option, "max").Raw

	case map[string]interface{}:
		minMaxOption.Min = getString(option["min"])
		minMaxOption.Max = getString(option["max"])
	}
	return minMaxOption
}
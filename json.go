package jdmapping

import (
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Json struct {
	oldJson string
	newJson string

	jdMapping *JDMapping
}

func NewJson(oldData, newData any, jdMapping *JDMapping) (j *Json, err error) {
	oldJson, err := marshalString(oldData)
	if err != nil {
		return
	}

	newJson, err := marshalString(newData)
	if err != nil {
		return
	}

	j = &Json{
		oldJson:   oldJson,
		newJson:   newJson,
		jdMapping: jdMapping,
	}

	err = j.handleArrayToMap()
	if err != nil {
		return
	}

	err = j.handleIsAddOP()
	if err != nil {
		return
	}

	err = j.handleExcludeKeys()
	if err != nil {
		return
	}

	return
}

func (j *Json) handleArrayToMap() (err error) {
	for _, item := range j.jdMapping.isArrayList {
		if err = j.convertArrayToMapFn(item.KeyList, item.ArrayIndexList, true); err != nil {
			return
		}
		if err = j.convertArrayToMapFn(item.KeyList, item.ArrayIndexList, false); err != nil {
			return
		}
	}
	return
}

func (j *Json) handleIsAddOP() (err error) {
	for _, item := range j.jdMapping.isAddOPList {
		mappingNode := findMappingNode(item.KeyList, j.jdMapping.config.MappingNode)
		if mappingNode == nil {
			continue
		}

		keyList := make([]string, len(item.KeyList))
		copy(keyList, item.KeyList)
		list := j.findJsonValueByKeyListFn(keyList, item.ArrayIndexList)

		for _, _item := range list {
			if item.MappingNode.IsAddOPFn(_item.oldValue, _item.newValue, mappingNode) {
				j.oldJson, err = sjson.Delete(j.oldJson, _item.Path)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (j *Json) handleExcludeKeys() (err error) {
	for _, item := range j.jdMapping.excludeKeysList {
		mappingNode := findMappingNode(item.KeyList, j.jdMapping.config.MappingNode)
		if mappingNode == nil {
			continue
		}

		keyList := make([]string, len(item.KeyList))
		copy(keyList, item.KeyList)

		list := j.findJsonValueByKeyListFnV2(keyList, item.ArrayIndexList, true)
		for _, _item := range list {
			for _, key := range mappingNode.ExcludeKeys {
				j.oldJson, err = sjson.Delete(j.oldJson, _item.Path+"."+key)
				if err != nil {
					return
				}
			}
		}

		list = j.findJsonValueByKeyListFnV2(keyList, item.ArrayIndexList, false)
		for _, _item := range list {
			for _, key := range mappingNode.ExcludeKeys {
				j.newJson, err = sjson.Delete(j.newJson, _item.Path+"."+key)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (j *Json) convertArrayToMapFn(keyList []string, arrayIndexList []int, isOld bool) (err error) {
	if len(arrayIndexList) == 0 {
		return
	}
	index := arrayIndexList[0]
	path := GetPath(keyList[:index+1])
	json := lo.Ternary(isOld, j.oldJson, j.newJson)

	result := gjson.Get(json, path)
	idxNumber := 0
	if result.IsArray() {
		newMap := make(map[string]any)
		for idx, item := range result.Array() {
			key := j.getConvertKey(keyList, idx, item)
			key = key + ":" + item.Get(key).String()
			if _, ok := newMap[key]; ok {
				idxNumber++
				key += "__" + cast.ToString(idxNumber)
			}

			newMap[key] = item.Value()
		}

		json, err = sjson.Set(json, path, newMap)
		if err != nil {
			return
		}

		if isOld {
			j.oldJson = json
		} else {
			j.newJson = json
		}
	}

	if len(arrayIndexList) != 0 && len(keyList) > index+1 {
		keys := gjson.Get(json, path+".@keys").Array()
		for _, item := range keys {
			keyList[index+1] = item.String()
			j.convertArrayToMapFn(keyList, arrayIndexList[1:], isOld)
		}
	}

	return
}

func (j *Json) getConvertKey(keyList []string, index int, value any) string {
	mappingNode := findMappingNode(keyList, j.jdMapping.config.MappingNode)
	if mappingNode == nil || (mappingNode.ConvertKey == "" && mappingNode.ConvertKeyFn == nil) {
		return cast.ToString(index)
	}

	if mappingNode.ConvertKey != "" {
		return mappingNode.ConvertKey
	}

	if text := mappingNode.ConvertKeyFn(value); text != "" {
		return text
	}

	return cast.ToString(index)
}

func (j *Json) findJsonValueByKeyListFn(keyList []string, arrayIndexList []int) (list []*JsonValue) {
	if len(arrayIndexList) == 0 {
		path := GetPath(keyList)
		oldValues := gjson.Get(j.oldJson, path)
		newValues := gjson.Get(j.newJson, path)
		return []*JsonValue{{Path: path, newValue: newValues.Value(), oldValue: oldValues.Value()}}
	}

	index := arrayIndexList[0]
	path := GetPath(keyList[:index+1])

	if len(keyList) > index+1 {
		keys := gjson.Get(j.newJson, path+".@keys").Array()
		for _, key := range keys {
			keyList[index+1] = key.String()
			list = append(list, j.findJsonValueByKeyListFn(keyList, arrayIndexList[1:])...)
		}
		return
	}

	oldValues := gjson.Get(j.oldJson, path)
	newValues := gjson.Get(j.newJson, path)
	return []*JsonValue{{Path: path, newValue: newValues.Value(), oldValue: oldValues.Value()}}
}

func (j *Json) findJsonValueByKeyListFnV2(keyList []string, arrayIndexList []int, isOld bool) (list []*JsonValueV2) {
	json := lo.Ternary(isOld, j.oldJson, j.newJson)

	if len(arrayIndexList) == 0 {
		path := GetPath(keyList)
		value := gjson.Get(json, path)
		return []*JsonValueV2{{Path: path, Value: value.Value()}}
	}

	index := arrayIndexList[0]
	path := GetPath(keyList[:index+1])

	if len(keyList) > index+1 {
		keys := gjson.Get(json, path+".@keys").Array()
		for _, key := range keys {
			keyList[index+1] = key.String()
			list = append(list, j.findJsonValueByKeyListFnV2(keyList, arrayIndexList[1:], isOld)...)
		}
		return
	}

	value := gjson.Get(json, path)
	return []*JsonValueV2{{Path: path, Value: value.Value()}}
}

func GetPath(keyList []string) string {
	if len(keyList) == 1 {
		return "@this"
	}
	return joinPoint(keyList[1:])
}

func findMappingNode(keyList []string, rootMappingNode *MappingNode) *MappingNode {
	if rootMappingNode == nil || len(keyList) <= 1 {
		return rootMappingNode
	}
	return findMappingNode(keyList[1:], rootMappingNode.Children[lo.Ternary(rootMappingNode.IsArray, ARRAY_NODE_INDEX_REPLACE_KEY, keyList[1])])
}

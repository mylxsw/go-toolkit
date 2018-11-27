package jsonutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JSONUtils a json utils object
type JSONUtils struct {
	message         []byte
	maxLevel        int
	obj             interface{}
	skipSimpleArray bool
}

// KvPair a kv pair
type KvPair struct {
	Key   string
	Value string
}

// New create a new json utils object and parse json to object
func New(message []byte, maxLevel int, skipSimpleArray bool) (*JSONUtils, error) {
	var obj interface{}
	if err := json.Unmarshal(message, &obj); err != nil {
		return nil, err
	}

	return &JSONUtils{
		message:         message,
		obj:             obj,
		maxLevel:        maxLevel,
		skipSimpleArray: skipSimpleArray,
	}, nil
}

// ToKvPairsArray convert to an array with all kv pair
func (ju *JSONUtils) ToKvPairsArray() []KvPair {
	return ju.createKvPairs(ju.obj, 1)
}

// ToKvPairs convert to a map with kv
func (ju *JSONUtils) ToKvPairs() map[string]string {
	kvPairs := make(map[string]string)
	for _, kv := range ju.ToKvPairsArray() {
		kvPairs[kv.Key] = kv.Value
	}

	return kvPairs
}

func (ju *JSONUtils) createKvPairs(obj interface{}, level int) []KvPair {
	kvPairs := make([]KvPair, 0)

	objValue := reflect.ValueOf(obj)
	if !objValue.IsValid() {
		return kvPairs
	}

	objType := objValue.Type()

	switch objType.Kind() {
	case reflect.Map:
		for _, key := range objValue.MapKeys() {
			keyStr := fmt.Sprintf("%s", key)
			value := objValue.MapIndex(key).Interface()

			subValues := ju.recursiveSubValue(keyStr, value, level+1)
			if len(subValues) == 0 {
				kvPairs = append(kvPairs, KvPair{
					Key:   keyStr,
					Value: "(null)",
				})
			}
			for _, kv := range subValues {
				kvPairs = append(kvPairs, kv)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < objValue.Len(); i++ {
			keyStr := fmt.Sprintf("[%d]", i)
			value := objValue.Index(i).Interface()

			subValues := ju.recursiveSubValue(keyStr, value, level+1)
			if len(subValues) == 0 {
				kvPairs = append(kvPairs, KvPair{
					Key:   keyStr,
					Value: "(null)",
				})
			}

			for _, kv := range subValues {
				kvPairs = append(kvPairs, kv)
			}
		}
	default:
	}

	return kvPairs
}

func (ju *JSONUtils) recursiveSubValue(keyStr string, value interface{}, level int) []KvPair {
	kvPairs := make([]KvPair, 0)
	reflectValue := reflect.ValueOf(value)
	if !reflectValue.IsValid() {
		return kvPairs
	}

	valueType := reflectValue.Type().Kind()

	switch valueType {
	case reflect.Slice, reflect.Map, reflect.Array:
		if ju.maxLevel > 0 && level > ju.maxLevel {
			valueJSON, _ := json.Marshal(value)
			kvPairs = append(kvPairs, KvPair{
				Key:   keyStr,
				Value: string(valueJSON),
			})

			break
		}

		// skip simple array
		if reflectValue.Len() > 0 && ju.skipSimpleArray && valueType != reflect.Map {
			subValue := reflect.ValueOf(reflectValue.Index(0).Interface())
			subValueKind := subValue.Kind()
			if subValueKind != reflect.Map && subValueKind != reflect.Slice && subValueKind != reflect.Array {
				valueJSON, _ := json.Marshal(value)
				kvPairs = append(kvPairs, KvPair{
					Key:   keyStr,
					Value: string(valueJSON),
				})

				break
			}
		}

		walkSub := ju.createKvPairs(value, level)
		// if sub value is empty, set value to [] or {}
		if len(walkSub) == 0 {
			emptyValue := "[]"
			if valueType == reflect.Map {
				emptyValue = "{}"
			}

			kvPairs = append(kvPairs, KvPair{
				Key:   keyStr,
				Value: emptyValue,
			})
		}

		for _, kv := range walkSub {
			kvPairs = append(kvPairs, KvPair{
				Key:   keyStr + "." + kv.Key,
				Value: kv.Value,
			})
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		kvPairs = append(kvPairs, KvPair{
			Key:   keyStr,
			Value: fmt.Sprintf("%d", value),
		})
	case reflect.Float32, reflect.Float64:
		kvPairs = append(kvPairs, KvPair{
			Key:   keyStr,
			Value: strconv.FormatFloat(value.(float64), 'f', -1, 64),
		})
	default:
		kvPairs = append(kvPairs, KvPair{
			Key:   keyStr,
			Value: fmt.Sprintf("%s", value),
		})
	}

	return kvPairs
}

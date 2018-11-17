package jsonutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// JSONUtils a json utils object
type JSONUtils struct {
	message []byte
	obj     interface{}
}

// KvPair a kv pair
type KvPair struct {
	Key   string
	Value string
}

// New create a new json utils object and parse json to object
func New(message []byte) (*JSONUtils, error) {
	var obj interface{}
	if err := json.Unmarshal(message, &obj); err != nil {
		return nil, err
	}

	return &JSONUtils{
		message: message,
		obj:     obj,
	}, nil
}

// ToKvPairsArray convert to an array with all kv pair
func (ju *JSONUtils) ToKvPairsArray() []KvPair {
	return ju.createKvPairs(ju.obj)
}

// ToKvPairs convert to a map with kv
func (ju *JSONUtils) ToKvPairs() map[string]string {
	kvPairs := make(map[string]string)
	for _, kv := range ju.ToKvPairsArray() {
		kvPairs[kv.Key] = kv.Value
	}

	return kvPairs
}

func (ju *JSONUtils) createKvPairs(obj interface{}) []KvPair {
	kvPairs := make([]KvPair, 0)

	objType := reflect.ValueOf(obj).Type()
	objValue := reflect.ValueOf(obj)

	switch objType.Kind() {
	case reflect.Map:
		for _, key := range objValue.MapKeys() {
			keyStr := fmt.Sprintf("%s", key)
			value := objValue.MapIndex(key).Interface()

			for _, kv := range ju.recursiveSubValue(keyStr, value) {
				kvPairs = append(kvPairs, kv)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < objValue.Len(); i++ {
			keyStr := fmt.Sprintf("[%d]", i)
			value := objValue.Index(i).Interface()

			for _, kv := range ju.recursiveSubValue(keyStr, value) {
				kvPairs = append(kvPairs, kv)
			}
		}
	default:
	}

	return kvPairs
}

func (ju *JSONUtils) recursiveSubValue(keyStr string, value interface{}) []KvPair {
	kvPairs := make([]KvPair, 0)
	valueType := reflect.ValueOf(value).Type().Kind()

	switch valueType {
	case reflect.Slice, reflect.Map, reflect.Array:
		for _, kv := range ju.createKvPairs(value) {
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

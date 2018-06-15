package collection

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrorInvalidDataType invalid data type error
	ErrorInvalidDataType = errors.New("Invalid data type")
)

// Collection is a data collection
type Collection struct {
	data     []interface{}
	dataType reflect.Type
}

// New create a new collection from data
func New(data interface{}) (*Collection, error) {
	dataKind := reflect.TypeOf(data).Kind()
	if dataKind != reflect.Array && dataKind != reflect.Slice {
		return nil, ErrorInvalidDataType
	}

	dataValue := reflect.ValueOf(data)
	dataArray := make([]interface{}, dataValue.Len())
	for i := 0; i < dataValue.Len(); i++ {
		dataArray[i] = dataValue.Index(i).Interface()
	}

	return &Collection{
		data: dataArray,
	}, nil
}

// MustNew create a new collection from data with error suppress
func MustNew(data interface{}) *Collection {
	res, err := New(data)
	if err != nil {
		panic(err.Error())
	}

	return res
}

// Filter iterates over elements of collection, return all element meet the needs
// filter(interface{}) bool
func (collection *Collection) Filter(filter interface{}) *Collection {
	if !IsFunction(filter, 1, 1) {
		panic("invalid callback function")
	}

	filterValue := reflect.ValueOf(filter)
	filterType := filterValue.Type()
	if filterType.Out(0).Kind() != reflect.Bool {
		panic("return argument should be a boolean")
	}

	results := make([]interface{}, 0)
	for _, item := range collection.data {
		if filterValue.Call([]reflect.Value{reflect.ValueOf(item)})[0].Interface().(bool) {
			results = append(results, item)
		}
	}

	return MustNew(results)
}

// Map manipulates an iteratee and transforms it to another type.
// mapFunc(interface{}) interface{}
func (collection *Collection) Map(mapFunc interface{}) *Collection {
	if !IsFunction(mapFunc, 1, 1) {
		panic("invalid callback function")
	}

	mapFuncValue := reflect.ValueOf(mapFunc)

	results := make([]interface{}, len(collection.data))
	for index, item := range collection.data {
		results[index] = mapFuncValue.Call([]reflect.Value{reflect.ValueOf(item)})[0].Interface()
	}

	return MustNew(results)
}

// Reduce Iteratively reduce the array to a single value using a callback function
// reduceFunc(carry interface{}, item interface{}) interface{}
func (collection *Collection) Reduce(reduceFunc interface{}, initial interface{}) interface{} {
	if !IsFunction(reduceFunc, 2, 1) {
		panic("invalid callback function")
	}

	reduceFuncValue := reflect.ValueOf(reduceFunc)

	previous := initial
	for _, item := range collection.data {
		previous = reduceFuncValue.Call([]reflect.Value{reflect.ValueOf(previous), reflect.ValueOf(item)})[0].Interface()
	}

	return previous
}

// All Get all of the items in the collection.
func (collection *Collection) All() []interface{} {
	return collection.data
}

// Each Execute a callback over each item.
func (collection *Collection) Each(eachFunc interface{}) {
	if !IsFunction(eachFunc) {
		panic("invalid callback function")
	}

	eachFuncValue := reflect.ValueOf(eachFunc)
	eachFuncType := eachFuncValue.Type()
	argumentNums := eachFuncType.NumIn()
	if argumentNums == 0 {
		panic("invalid callback function")
	}

	for index, item := range collection.data {
		if argumentNums == 1 {
			eachFuncValue.Call([]reflect.Value{reflect.ValueOf(item)})
		} else {
			eachFuncValue.Call([]reflect.Value{reflect.ValueOf(item), reflect.ValueOf(index)})
		}
	}
}

// Append append new items to collection
func (collection *Collection) Append(items ...interface{}) {
	collection.data = append(collection.data, items...)
}

// IsEmpty Determine if the collection is empty or not.
func (collection *Collection) IsEmpty() bool {
	return len(collection.data) == 0
}

// ToString print the data element
func (collection *Collection) ToString() string {
	return fmt.Sprint(collection.data)
}

// IsFunction returns if the argument is a function.
func IsFunction(in interface{}, num ...int) bool {
	funcType := reflect.TypeOf(in)

	result := funcType.Kind() == reflect.Func

	if len(num) >= 1 {
		result = result && funcType.NumIn() == num[0]
	}

	if len(num) == 2 {
		result = result && funcType.NumOut() == num[1]
	}

	return result
}

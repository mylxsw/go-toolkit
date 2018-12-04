package container

import (
	"errors"
	"reflect"
	"sync"
)

var (
	// ErrObjectNotFound is an error object represent object not found
	ErrObjectNotFound = errors.New("the object can not be found in container")
	// ErrArgNotInstanced is an erorr object represent arg not instanced
	ErrArgNotInstanced = errors.New("the arg can not be found in container")
	// ErrInvalidReturnValueCount is an error object represent return values count not match
	ErrInvalidReturnValueCount = errors.New("invalid return value count")
	// ErrRepeatedBind is an error object represent bind a value repeated
	ErrRepeatedBind = errors.New("can not bind a value with repeated key")
	// ErrInvalidArgs is an error object represent invalid args
	ErrInvalidArgs = errors.New("invalid args")
)

// Entity represent a entity in container
type Entity struct {
	lock sync.RWMutex

	key            interface{} // entity key
	initializeFunc interface{} // initializeFunc is a func to initialize entity
	value          interface{}
	typ            reflect.Type
	index          int // the index in the container

	prototype bool
	c         *Container
}

// Value instance value if not initiailzed
func (e *Entity) Value() (interface{}, error) {
	if e.prototype {
		return e.createValue()
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if e.value == nil {
		val, err := e.createValue()
		if err != nil {
			return nil, err
		}

		e.value = val
	}

	return e.value, nil
}

func (e *Entity) createValue() (interface{}, error) {
	initializeValue := reflect.ValueOf(e.initializeFunc)
	argValues, err := e.c.funcArgs(initializeValue.Type())
	if err != nil {
		return nil, err
	}

	returnValues := reflect.ValueOf(e.initializeFunc).Call(argValues)
	if len(returnValues) != 2 {
		return nil, ErrInvalidReturnValueCount
	}

	if !returnValues[1].IsNil() && returnValues[1].Interface() != nil {
		return nil, returnValues[1].Interface().(error)
	}

	return returnValues[0].Interface(), nil
}

// Container is a dependency injection container
type Container struct {
	lock sync.RWMutex

	objects      map[interface{}]*Entity
	objectSlices []*Entity
}

// New create a new container
func New() *Container {
	return &Container{
		objects:      make(map[interface{}]*Entity),
		objectSlices: make([]*Entity, 0),
	}
}

// Prototype bind a prototype
// initialize func(...) (value, error)
func (c *Container) Prototype(initialize interface{}) error {
	return c.Bind(initialize, true)
}

// PrototypeWithKey bind a prototype with key
// initialize func(...) (value, error)
func (c *Container) PrototypeWithKey(key interface{}, initialize interface{}) error {
	return c.BindWithKey(key, initialize, true)
}

// Singleton bind a singleton
// initialize func(...) (value, error)
func (c *Container) Singleton(initialize interface{}) error {
	return c.Bind(initialize, false)
}

// SingletonWithKey bind a singleton with key
// initialize func(...) (value, error)
func (c *Container) SingletonWithKey(key interface{}, initialize interface{}) error {
	return c.BindWithKey(key, initialize, false)
}

// BindValue bing a value to container
func (c *Container) BindValue(key interface{}, value interface{}) error {
	if value == nil {
		return ErrInvalidArgs
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.objects[key]; ok {
		return ErrRepeatedBind
	}

	entity := Entity{
		initializeFunc: nil,
		key:            key,
		typ:            reflect.TypeOf(value),
		value:          value,
		index:          len(c.objectSlices),
		c:              c,
		prototype:      false,
	}

	c.objects[key] = &entity
	c.objectSlices = append(c.objectSlices, &entity)

	return nil
}

// Bind bind a initialize for object
// initialize func(...) (value, error)
func (c *Container) Bind(initialize interface{}, prototype bool) error {
	if !reflect.ValueOf(initialize).IsValid() {
		return ErrInvalidArgs
	}

	initializeType := reflect.ValueOf(initialize).Type()
	if initializeType.NumOut() != 2 || !c.isErrorType(initializeType.Out(1)) {
		return ErrInvalidArgs
	}

	typ := initializeType.Out(0)
	return c.bindWith(typ, typ, initialize, prototype)
}

// BindWithKey bind a initialize for object with a key
// initialize func(...) (value, error)
func (c *Container) BindWithKey(key interface{}, initialize interface{}, prototype bool) error {
	if !reflect.ValueOf(initialize).IsValid() {
		return ErrInvalidArgs
	}

	initializeType := reflect.ValueOf(initialize).Type()
	if initializeType.NumOut() != 2 || !c.isErrorType(initializeType.Out(1)) {
		return ErrInvalidArgs
	}

	return c.bindWith(key, initializeType.Out(0), initialize, prototype)
}

// Resolve inject args for func by callback
// callback func(...)
func (c *Container) Resolve(callback interface{}) error {
	callbackValue := reflect.ValueOf(callback)
	if !callbackValue.IsValid() {
		return ErrInvalidArgs
	}

	args, err := c.funcArgs(callbackValue.Type())
	if err != nil {
		return err
	}

	callbackValue.Call(args)
	return nil
}

// Get get instance by key from container
func (c *Container) Get(key interface{}) (interface{}, error) {
	keyReflectType, ok := key.(reflect.Type)
	if !ok {
		keyReflectType = reflect.TypeOf(key)
	}

	for _, obj := range c.objectSlices {

		if obj.key == key || obj.key == keyReflectType {
			return obj.Value()
		}

		if obj.typ.AssignableTo(keyReflectType) {
			return obj.Value()
		}
	}

	return nil, ErrObjectNotFound
}

func (c *Container) bindWith(key interface{}, typ reflect.Type, initialize interface{}, prototype bool) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.objects[key]; ok {
		return ErrRepeatedBind
	}

	entity := Entity{
		initializeFunc: initialize,
		key:            key,
		typ:            typ,
		value:          nil,
		index:          len(c.objectSlices),
		c:              c,
		prototype:      prototype,
	}

	c.objects[key] = &entity
	c.objectSlices = append(c.objectSlices, &entity)

	return nil
}

func (c *Container) isErrorType(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

func (c *Container) funcArgs(t reflect.Type) ([]reflect.Value, error) {
	argsSize := t.NumIn()
	argValues := make([]reflect.Value, argsSize)
	for i := 0; i < argsSize; i++ {
		argType := t.In(i)
		val, err := c.instanceOfType(argType)
		if err != nil {
			return argValues, err
		}

		argValues[i] = val
	}

	return argValues, nil
}

func (c *Container) instanceOfType(t reflect.Type) (reflect.Value, error) {
	if reflect.TypeOf(c).AssignableTo(t) {
		return reflect.ValueOf(c), nil
	}

	arg, err := c.Get(t)
	if err != nil {
		return reflect.Value{}, ErrArgNotInstanced
	}

	return reflect.ValueOf(arg), nil
}

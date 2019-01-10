package events

import (
	"fmt"
	"reflect"
	"sync"
)

// Listener is a event listener
type Listener interface{}

// EventStore is a interface for event store
type EventStore interface {
	Listen(eventName string, listener Listener)
	Publish(eventName string, evt interface{})
	SetManager(*EventManager)
}

// EventManager is a manager for event dispatch
type EventManager struct {
	store EventStore
	lock  sync.RWMutex
}

// NewEventManager create a eventManager
func NewEventManager(store EventStore) *EventManager {
	manager := &EventManager{
		store: store,
	}

	store.SetManager(manager)

	return manager
}

// Listen create a relation from event to listners
func (em *EventManager) Listen(listeners ...Listener) {
	em.lock.Lock()
	defer em.lock.Unlock()

	for _, listener := range listeners {
		listenerType := reflect.TypeOf(listener)
		if listenerType.Kind() != reflect.Func {
			panic("listener must be a function")
		}

		if listenerType.NumIn() != 1 {
			panic("listener must be a function with only one arguemnt")
		}

		if listenerType.In(0).Kind() != reflect.Struct {
			panic("listener must be a function with only on argument of type struct")
		}

		em.store.Listen(fmt.Sprintf("%s", listenerType.In(0)), listener)
	}
}

// Publish a event
func (em *EventManager) Publish(evt interface{}) {
	em.lock.RLock()
	defer em.lock.RUnlock()

	em.store.Publish(fmt.Sprintf("%s", reflect.TypeOf(evt)), evt)
}

// Call trigger listener to execute
func (em *EventManager) Call(evt interface{}, listener Listener) {
	reflect.ValueOf(listener).Call([]reflect.Value{reflect.ValueOf(evt)})
}

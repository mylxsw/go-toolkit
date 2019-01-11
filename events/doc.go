/*
Package events 实现了应用中的事件机制，用于业务逻辑解耦。

例如

	eventManager := events.NewEventManager(events.NewMemoryEventStore(false))
	eventManager.Listen(func(evt UserCreatedEvent) {
		t.Logf("user created: id=%s, name=%s", evt.ID, evt.UserName)

		if evt.ID != "111" {
			t.Error("test failed")
		}
	})

	eventManager.Publish(UserCreatedEvent{
		ID:       "111",
		UserName: "李逍遥",
	})

*/
package events

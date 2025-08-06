package eventBus

import "sync"

type Event interface {
	Name() string
}

type UserLoggedInEvent struct {
	Phone string
}

func (u UserLoggedInEvent) Name() string {
	return "UserLoggedInEvent"
}

func (u UserLoggedInEvent) GetPhone() string {
	return u.Phone
}

type EventBus struct {
	channels map[string]chan Event
	wg       sync.WaitGroup
}

func NewEventBus() *EventBus {
	return &EventBus{
		channels: make(map[string]chan Event),
	}
}

func (bus *EventBus) Publish(event Event) {
	eventName := event.Name()

	if _, ok := bus.channels[eventName]; ok {
		bus.channels[eventName] <- event
	}
}

func (bus *EventBus) Subscribe(event Event) {
	eventName := event.Name()

	if _, ok := bus.channels[eventName]; !ok {
		bus.channels[eventName] = make(chan Event, 10)
	}

	bus.wg.Add(1)
	go bus.processEvents(eventName)
}

func (bus *EventBus) processEvents(eventName string) {
	defer bus.wg.Done()
	_ = bus.channels[eventName]
}

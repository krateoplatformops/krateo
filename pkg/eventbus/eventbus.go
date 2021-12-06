// Package eventbus provides event publisher/subscriber support.
// It's inspired by https://github.com/asaskevich/EventBus
// There are differences in API and implementation
//  - deadlock prevention: this eventbus doesn't lock mutex when callbacks are called
//  - API based on EventID and Event interface instead of strings and variadic arguments list

package eventbus

import (
	"sync"
)

// EventID identifies events topic.
type EventID string

// Event must be implemented by anything that can be published
type Event interface {
	EventID() EventID
}

// EventHandler is function that can be subscribed to the event
type EventHandler func(event Event)

// Subscription represents active event subscription
type Subscription struct {
	eventID EventID
	id      uint64
}

// BusSubscriber allows to subscribe/unsubscribe own event handlers
type BusSubscriber interface {
	Subscribe(eventID EventID, cb EventHandler) Subscription
	Unsubscribe(id Subscription)
}

// BusPublisher allows to publish own events
type BusPublisher interface {
	Publish(event Event)
}

// Bus allows to subscribe/unsubscribe to external events and publish own events
type Bus interface {
	BusSubscriber
	BusPublisher
}

// New returns new event bus
func New() Bus {
	b := &bus{
		infos: make(map[EventID]subscriptionInfoList),
	}
	return b
}

type subscriptionInfo struct {
	id uint64
	cb EventHandler
}

type subscriptionInfoList []*subscriptionInfo

type bus struct {
	lock   sync.Mutex
	nextID uint64
	infos  map[EventID]subscriptionInfoList
}

func (bus *bus) Subscribe(eventID EventID, cb EventHandler) Subscription {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	id := bus.nextID
	bus.nextID++
	sub := &subscriptionInfo{
		id: id,
		cb: cb,
	}
	bus.infos[eventID] = append(bus.infos[eventID], sub)
	return Subscription{
		eventID: eventID,
		id:      id,
	}
}

func (bus *bus) Unsubscribe(subscription Subscription) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if infos, ok := bus.infos[subscription.eventID]; ok {
		for idx, info := range infos {
			if info.id == subscription.id {
				infos = append(infos[:idx], infos[idx+1:]...)
				break
			}
		}
		if len(infos) == 0 {
			delete(bus.infos, subscription.eventID)
		} else {
			bus.infos[subscription.eventID] = infos
		}
	}
}

func (bus *bus) Publish(event Event) {
	infos := bus.copySubscriptions(event.EventID())
	for _, sub := range infos {
		sub.cb(event)
	}
}

func (bus *bus) copySubscriptions(eventID EventID) subscriptionInfoList {
	// External code may subscribe/unsubscribe during iteration over callbacks,
	//  so we need to copy subscribers to invoke callbacks.

	bus.lock.Lock()
	defer bus.lock.Unlock()
	if infos, ok := bus.infos[eventID]; ok {
		return infos
	}
	return subscriptionInfoList{}
}

package eventbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	eventSolarEclipse = EventID("solar_eclipse")
	eventMoonEclipse  = EventID("moon_eclipse")
)

type solarEclipseEvent struct {
	duration time.Duration
}

func (e *solarEclipseEvent) EventID() EventID {
	return eventSolarEclipse
}

type moonEclipseEvent struct {
	duration time.Duration
}

func (e *moonEclipseEvent) EventID() EventID {
	return eventMoonEclipse
}

func TestBus_SubscribePublish(t *testing.T) {
	bus := New()
	hadEvent := false
	duration := 100 * time.Second

	bus.Subscribe(eventSolarEclipse, func(e Event) {
		assert.Equal(t, e.EventID(), eventSolarEclipse)
		se := e.(*solarEclipseEvent)
		assert.Equal(t, se.duration, duration)
		hadEvent = true
	})
	bus.Subscribe(eventMoonEclipse, func(e Event) {
		t.Fatalf("should never be called")
	})
	assert.Equal(t, hadEvent, false)

	bus.Publish(&solarEclipseEvent{
		duration: duration,
	})
	assert.Equal(t, hadEvent, true)
}

func TestBus_PublishIncompatibleEvent(t *testing.T) {
	bus := New()
	duration := 100 * time.Second
	bus.Subscribe(eventMoonEclipse, func(e Event) {
		t.Fatalf("should never be called")
	})

	bus.Publish(&solarEclipseEvent{
		duration: duration,
	})
}

func TestBus_SubscribeUnsubscribe(t *testing.T) {
	bus := New()
	hadEvent := false
	duration := 42 * time.Millisecond

	id := bus.Subscribe(eventMoonEclipse, func(e Event) {
		assert.Equal(t, e.EventID(), eventMoonEclipse)
		se := e.(*moonEclipseEvent)
		assert.Equal(t, se.duration, duration)
		hadEvent = true
	})
	bus.Publish(&moonEclipseEvent{
		duration: duration,
	})
	assert.Equal(t, hadEvent, true)

	hadEvent = false
	bus.Unsubscribe(id)
	bus.Publish(&moonEclipseEvent{
		duration: duration,
	})
	assert.Equal(t, hadEvent, false)
}

func TestBus_SubscribeMultiple(t *testing.T) {
	moonEventCount := 0
	moonEclipseDuration := 16 * time.Second
	onMoonEclipse := func(e Event) {
		assert.Equal(t, e.EventID(), eventMoonEclipse)
		moonEventCount++
		se := e.(*moonEclipseEvent)
		assert.Equal(t, se.duration, moonEclipseDuration)
	}

	solarEventCount := 0
	solarEclipseDuration := 77 * time.Millisecond
	onSolarEclipse := func(e Event) {
		assert.Equal(t, e.EventID(), eventSolarEclipse)
		solarEventCount++
		se := e.(*solarEclipseEvent)
		assert.Equal(t, se.duration, solarEclipseDuration)
	}

	bus := New()
	publishMoon := func() {
		bus.Publish(&moonEclipseEvent{
			duration: moonEclipseDuration,
		})
	}
	publishSolar := func() {
		bus.Publish(&solarEclipseEvent{
			duration: solarEclipseDuration,
		})
	}

	id1 := bus.Subscribe(eventMoonEclipse, onMoonEclipse)
	id2 := bus.Subscribe(eventSolarEclipse, onSolarEclipse)
	id3 := bus.Subscribe(eventMoonEclipse, onMoonEclipse)

	publishMoon()
	assert.Equal(t, moonEventCount, 2)
	assert.Equal(t, solarEventCount, 0)

	publishSolar()
	assert.Equal(t, moonEventCount, 2)
	assert.Equal(t, solarEventCount, 1)

	bus.Unsubscribe(id1)
	bus.Unsubscribe(id2)
	publishMoon()
	publishSolar()
	assert.Equal(t, moonEventCount, 3)
	assert.Equal(t, solarEventCount, 1)

	bus.Unsubscribe(id3)
	publishMoon()
	publishSolar()
	assert.Equal(t, moonEventCount, 3)
	assert.Equal(t, solarEventCount, 1)
}
func TestBus_PublishRecursive(t *testing.T) {
	moonEventCount := 0

	bus := New()
	publishMoon := func(duration time.Duration) {
		bus.Publish(&moonEclipseEvent{
			duration: duration,
		})
	}

	onMoonEclipse := func(e Event) {
		assert.Equal(t, e.EventID(), eventMoonEclipse)
		moonEventCount++
		se := e.(*moonEclipseEvent)

		if se.duration < 16*time.Second {
			publishMoon(2 * se.duration)
		}
	}

	bus.Subscribe(eventMoonEclipse, onMoonEclipse)
	publishMoon(1 * time.Second)

	assert.Equal(t, moonEventCount, 5)
}

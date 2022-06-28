package events

import (
	"fmt"

	"github.com/krateoplatformops/krateo/pkg/eventbus"
)

const (
	StartWaitEventID = eventbus.EventID("log:startWait")
	StopWaitEventID  = eventbus.EventID("log:stopWait")
	DoneEventID      = eventbus.EventID("log:done")
	DebugEventID     = eventbus.EventID("log:debug")
	WarningEventID   = eventbus.EventID("log:warn")
)

func NewStartWaitEvent(s string, args ...interface{}) *StartWaitEvent {
	return &StartWaitEvent{fmt.Sprintf(s, args...)}
}

type StartWaitEvent struct {
	message string
}

func (e *StartWaitEvent) EventID() eventbus.EventID {
	return StartWaitEventID
}

func (e *StartWaitEvent) Message() string {
	return e.message
}

func NewStopWaitEvent() *StopWaitEvent {
	return &StopWaitEvent{}
}

type StopWaitEvent struct{}

func (e *StopWaitEvent) EventID() eventbus.EventID {
	return StopWaitEventID
}

func NewDoneEvent(s string, args ...interface{}) *DoneEvent {
	return &DoneEvent{fmt.Sprintf(s, args...)}
}

type DoneEvent struct {
	message string
}

func (e *DoneEvent) EventID() eventbus.EventID {
	return DoneEventID
}

func (e *DoneEvent) Message() string {
	return e.message
}

func NewDebugEvent(s string, args ...interface{}) *DebugEvent {
	return &DebugEvent{fmt.Sprintf(s, args...)}
}

type DebugEvent struct {
	message string
}

func (e *DebugEvent) EventID() eventbus.EventID {
	return DebugEventID
}

func (e *DebugEvent) Message() string {
	return e.message
}

func NewWarningEvent(s string, args ...interface{}) *WarningEvent {
	return &WarningEvent{fmt.Sprintf(s, args...)}
}

type WarningEvent struct {
	message string
}

func (e *WarningEvent) EventID() eventbus.EventID {
	return WarningEventID
}

func (e *WarningEvent) Message() string {
	return e.message
}

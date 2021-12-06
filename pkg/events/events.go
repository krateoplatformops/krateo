package events

import (
	"github.com/krateoplatformops/krateoctl/pkg/eventbus"
)

const (
	StartWaitEventID = eventbus.EventID("log:startWait")
	StopWaitEventID  = eventbus.EventID("log:stopWait")
	DoneEventID      = eventbus.EventID("log:done")
	DebugEventID     = eventbus.EventID("log:debug")
)

func NewStartWaitEvent(s string) *StartWaitEvent {
	return &StartWaitEvent{s}
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

func NewDoneEvent(s string) *DoneEvent {
	return &DoneEvent{s}
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

func NewDebugEvent(s string) *DebugEvent {
	return &DebugEvent{s}
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

package cmd

import (
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
)

func updateLog(l log.Logger) eventbus.EventHandler {
	return func(e eventbus.Event) {
		switch e.EventID() {
		case events.DebugEventID:
			evt := e.(*events.DebugEvent)
			l.Debug(evt.Message())

		case events.StartWaitEventID:
			evt := e.(*events.StartWaitEvent)
			l.StartWait(evt.Message())

		case events.StopWaitEventID:
			l.StopWait()

		case events.DoneEventID:
			evt := e.(*events.DoneEvent)
			l.StopWait()
			l.Done(evt.Message())
		}
	}
}

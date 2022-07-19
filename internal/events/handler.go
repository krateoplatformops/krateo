package events

import (
	"github.com/krateoplatformops/krateo/internal/eventbus"
	"github.com/krateoplatformops/krateo/internal/log"
)

func LogHandler(l log.Logger) eventbus.EventHandler {
	return func(e eventbus.Event) {
		switch e.EventID() {
		case DebugEventID:
			evt := e.(*DebugEvent)
			l.Debug(evt.Message())

		case WarningEventID:
			evt := e.(*WarningEvent)
			l.Warn(evt.Message())

		case StartWaitEventID:
			evt := e.(*StartWaitEvent)
			l.StartWait(evt.Message())

		case StopWaitEventID:
			l.StopWait()

		case DoneEventID:
			evt := e.(*DoneEvent)
			l.StopWait()
			l.Done(evt.Message())
		}
	}
}

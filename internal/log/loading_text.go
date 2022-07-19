package log

import (
	"fmt"
	"io"
	"time"

	"github.com/mgutz/ansi"
	"golang.org/x/term"
)

const (
	waitInterval = time.Millisecond * 150
)

var (
	spinnerFrames = []string{
		"⠈⠁", "⠈⠑", "⠈⠱", "⠈⡱", "⢀⡱", "⢄⡱", "⢄⡱", "⢆⡱", "⢎⡱", "⢎⡰", "⢎⡠", "⢎⡀", "⢎⠁", "⠎⠁", "⠊⠁",
	}
)

type loadingText struct {
	Stream         io.Writer
	Message        string
	StartTimestamp int64

	loadingRune int
	isShown     bool
	stopChan    chan bool
}

func (l *loadingText) Start() {
	l.isShown = false
	l.StartTimestamp = time.Now().UnixNano()

	if l.stopChan == nil {
		l.stopChan = make(chan bool)
	}

	go func() {
		l.render()

		for {
			select {
			case <-l.stopChan:
				return
			case <-time.After(waitInterval):
				l.render()
			}
		}
	}()
}

func (l *loadingText) getLoadingChar() string {
	max := len(spinnerFrames)
	idx := l.loadingRune % max
	loadingChar := spinnerFrames[idx]

	l.loadingRune++
	if l.loadingRune > max {
		l.loadingRune = 0
	}

	return loadingChar
}

//nolint:errcheck
func (l *loadingText) render() {
	if !l.isShown {
		l.isShown = true
	} else {
		l.Stream.Write([]byte("\r"))
	}
	messagePrefix := []byte("[wait] ")

	l.Stream.Write([]byte(ansi.Color(string(messagePrefix), "cyan+b")))

	timeElapsed := fmt.Sprintf("%d", (time.Now().UnixNano()-l.StartTimestamp)/int64(time.Second))
	message := []byte(l.getLoadingChar() + " " + l.Message)
	messageSuffix := " (" + timeElapsed + "s)"
	prefixLength := len(messagePrefix)
	suffixLength := len(messageSuffix)

	tWidth, _, err := term.GetSize(0)
	if err == nil && prefixLength+len(message)+suffixLength > tWidth {
		dots := []byte("...")

		maxMessageLength := tWidth - (prefixLength + suffixLength + len(dots))
		if maxMessageLength > 0 {
			message = append(message[:maxMessageLength], dots...)
		}
	}

	message = append(message, messageSuffix...)
	l.Stream.Write(message)
}

//nolint:errcheck
func (l *loadingText) Stop() {
	l.stopChan <- true
	l.Stream.Write([]byte("\r"))

	messageLength := len(l.Message) + 20

	for i := 0; i < messageLength; i++ {
		l.Stream.Write([]byte(" "))
	}

	l.Stream.Write([]byte("\r"))
}

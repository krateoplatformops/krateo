package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	goansi "github.com/k0kubun/go-ansi"
	"github.com/mgutz/ansi"
)

var (
	stdout = goansi.NewAnsiStdout()
	stderr = goansi.NewAnsiStderr()
)

type stdoutLogger struct {
	logMutex sync.Mutex
	level    Level

	loadingText *loadingText

	//	fileLogger Logger
}

type fnTypeInformation struct {
	tag      string
	color    string
	logLevel Level
	stream   io.Writer
}

var fnTypeInformationMap = map[logFunctionType]*fnTypeInformation{
	debugFn: {
		tag:      "[debug]  ",
		color:    "yellow",
		logLevel: DebugLevel,
		stream:   stdout,
	},
	infoFn: {
		tag:      "[info]   ",
		color:    "cyan+b",
		logLevel: InfoLevel,
		stream:   stdout,
	},
	warnFn: {
		tag:      "[warn]   ",
		color:    "red+b",
		logLevel: WarnLevel,
		stream:   stdout,
	},
	errorFn: {
		tag:      "[error]  ",
		color:    "red+b",
		logLevel: ErrorLevel,
		stream:   stdout,
	},
	fatalFn: {
		tag:      "[fatal]  ",
		color:    "red+b",
		logLevel: FatalLevel,
		stream:   stdout,
	},
	panicFn: {
		tag:      "[panic]  ",
		color:    "red+b",
		logLevel: PanicLevel,
		stream:   stderr,
	},
	doneFn: {
		tag:      "[done] √ ",
		color:    "green+b",
		logLevel: InfoLevel,
		stream:   stdout,
	},
	failFn: {
		tag:      "[fail] X ",
		color:    "red+b",
		logLevel: ErrorLevel,
		stream:   stdout,
	},
}

//nolint:errcheck
func (s *stdoutLogger) writeMessage(fnType logFunctionType, message string) {
	fnInformation := fnTypeInformationMap[fnType]
	if s.level >= fnInformation.logLevel {
		if s.loadingText != nil {
			s.loadingText.Stop()
		}

		sc := bufio.NewScanner(strings.NewReader(message))
		for sc.Scan() {
			fnInformation.stream.Write([]byte(ansi.Color(fnInformation.tag, fnInformation.color)))
			fnInformation.stream.Write([]byte(fmt.Sprintln(sc.Text())))
		}

		if s.loadingText != nil && fnType != fatalFn {
			s.loadingText.Start()
		}
	}
}

// StartWait prints a wait message until StopWait is called
func (s *stdoutLogger) StartWait(message string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		if s.loadingText.Message == message {
			return
		}

		s.loadingText.Stop()
		s.loadingText = nil
	}

	if s.level >= InfoLevel {
		s.loadingText = &loadingText{
			Message: message,
			Stream:  goansi.NewAnsiStdout(),
		}

		s.loadingText.Start()
	}
}

// StartWait prints a wait message until StopWait is called
func (s *stdoutLogger) StopWait() {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		s.loadingText.Stop()
		s.loadingText = nil
	}
}

func (s *stdoutLogger) Debug(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(debugFn, fmt.Sprintln(args...))
}

func (s *stdoutLogger) Debugf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(debugFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Info(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(infoFn, fmt.Sprintln(args...))
}

func (s *stdoutLogger) Infof(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(infoFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Warn(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(warnFn, fmt.Sprintln(args...))
}

func (s *stdoutLogger) Warnf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(warnFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Error(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(errorFn, fmt.Sprintln(args...))
}

func (s *stdoutLogger) Errorf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(errorFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Fatal(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	msg := fmt.Sprintln(args...)

	s.writeMessage(fatalFn, msg)

	os.Exit(1)
}

func (s *stdoutLogger) Fatalf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	msg := fmt.Sprintf(format, args...)

	s.writeMessage(fatalFn, msg+"\n")

	os.Exit(1)
}

func (s *stdoutLogger) Panic(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(panicFn, fmt.Sprintln(args...))

	panic(fmt.Sprintln(args...))
}

func (s *stdoutLogger) Panicf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(panicFn, fmt.Sprintf(format, args...)+"\n")

	panic(fmt.Sprintf(format, args...))
}

func (s *stdoutLogger) Done(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(doneFn, fmt.Sprintln(args...))

}

func (s *stdoutLogger) Donef(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(doneFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Fail(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(failFn, fmt.Sprintln(args...))
}

func (s *stdoutLogger) Failf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(failFn, fmt.Sprintf(format, args...)+"\n")
}

func (s *stdoutLogger) Print(level Level, args ...interface{}) {
	switch level {
	case InfoLevel:
		s.Info(args...)
	case DebugLevel:
		s.Debug(args...)
	case WarnLevel:
		s.Warn(args...)
	case ErrorLevel:
		s.Error(args...)
	case PanicLevel:
		s.Panic(args...)
	case FatalLevel:
		s.Fatal(args...)
	}
}

func (s *stdoutLogger) Printf(level Level, format string, args ...interface{}) {
	switch level {
	case InfoLevel:
		s.Infof(format, args...)
	case DebugLevel:
		s.Debugf(format, args...)
	case WarnLevel:
		s.Warnf(format, args...)
	case ErrorLevel:
		s.Errorf(format, args...)
	case PanicLevel:
		s.Panicf(format, args...)
	case FatalLevel:
		s.Fatalf(format, args...)
	}
}

func (s *stdoutLogger) SetLevel(level Level) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.level = level
}

func (s *stdoutLogger) GetLevel() Level {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	return s.level
}

func (s *stdoutLogger) Write(message []byte) (int, error) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.level >= InfoLevel {
		if s.loadingText != nil {
			s.loadingText.Stop()
		}

		n, err := fnTypeInformationMap[infoFn].stream.Write(message)

		if s.loadingText != nil {
			s.loadingText.Start()
		}

		return n, err
	}

	return len(message), nil
}

//nolint:errcheck
func (s *stdoutLogger) WriteString(message string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.level >= InfoLevel {
		if s.loadingText != nil {
			s.loadingText.Stop()
		}

		fnTypeInformationMap[infoFn].stream.Write([]byte(message))

		if s.loadingText != nil {
			s.loadingText.Start()
		}
	}
}

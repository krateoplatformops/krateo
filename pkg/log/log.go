//nolint:errcheck
package log

import (
	"strings"

	"github.com/mgutz/ansi"
)

var defaultLog Logger = &stdoutLogger{
	level: DebugLevel,
}

// GetInstance returns the Logger instance
func GetInstance() Logger {
	return defaultLog
}

// SetInstance sets the default logger instance
func SetInstance(logger Logger) {
	defaultLog = logger
}

// WriteColored writes a message in color
func writeColored(message string, color string) {
	defaultLog.Write([]byte(ansi.Color(message, color)))
}

//SetFakePrintTable is a testing tool that allows overwriting the function PrintTable
func SetFakePrintTable(fake func(s Logger, header []string, values [][]string)) {
	fakePrintTable = fake
}

var fakePrintTable func(s Logger, header []string, values [][]string)

// PrintTable prints a table with header columns and string values
func PrintTable(s Logger, header []string, values [][]string) {
	if fakePrintTable != nil {
		fakePrintTable(s, header, values)
		return
	}

	columnLengths := make([]int, len(header))

	for k, v := range header {
		columnLengths[k] = len(v)
	}

	// Get maximum column length
	for _, v := range values {
		for key, value := range v {
			if len(value) > columnLengths[key] {
				columnLengths[key] = len(value)
			}
		}
	}

	s.Write([]byte("\n"))

	// Print Header
	for key, value := range header {
		writeColored(" "+value+"  ", "green+b")

		padding := columnLengths[key] - len(value)

		if padding > 0 {
			s.Write([]byte(strings.Repeat(" ", padding)))
		}
	}

	s.Write([]byte("\n"))

	if len(values) == 0 {
		s.Write([]byte(" No entries found\n"))
	}

	// Print Values
	for _, v := range values {
		for key, value := range v {
			s.Write([]byte(" " + value + "  "))

			padding := columnLengths[key] - len(value)

			if padding > 0 {
				s.Write([]byte(strings.Repeat(" ", padding)))
			}
		}

		s.Write([]byte("\n"))
	}

	s.Write([]byte("\n"))
}

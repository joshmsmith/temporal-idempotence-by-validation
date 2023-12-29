package utils

import (
	"fmt"
	golog "log"
	"os"
	"strings"

	"go.temporal.io/sdk/log"
)

/* TClientLogger - Example custom client Logger implementation on top of standard log.Logger. */
type TClientLogger struct {
	logger        *golog.Logger
	globalKeyvals string
}

/* NewTClientLogger creates new instance of TClientLogger. */
func NewTClientLogger() *TClientLogger {
	return &TClientLogger{logger: golog.New(os.Stdout, "", golog.LstdFlags)}
}

/* println */
func (l *TClientLogger) println(level, msg string, keyvals []interface{}) {
	// To avoid extra space when globalKeyvals is not specified.
	if l.globalKeyvals == "" {
		l.logger.Println(append([]interface{}{level, msg}, keyvals...)...)
	} else {
		//l.logger.Println(append([]interface{}{level, msg, l.globalKeyvals}, keyvals...)...)
		l.logger.Println(append([]interface{}{level, msg}, keyvals...)...)
	}
}

/* Debug writes message to the log. */
func (l *TClientLogger) Debug(msg string, keyvals ...interface{}) {
	if log_level == "debug" {
		l.println("DEBUG", msg, keyvals)
	}
}

/* Info writes message to the log. */
func (l *TClientLogger) Info(msg string, keyvals ...interface{}) {
	l.println("INFO ", msg, keyvals)
}

/* QInfo writes message to the log. */
func (l *TClientLogger) QInfo(msg string, keyvals ...interface{}) {
	l.println("", msg, keyvals)
}

/* Warn writes message to the log. */
func (l *TClientLogger) Warn(msg string, keyvals ...interface{}) {
	l.println("WARN ", msg, keyvals)
}

/* Error writes message to the log. */
func (l *TClientLogger) Error(msg string, keyvals ...interface{}) {
	l.println("ERROR", msg, keyvals)
}

/* With returns new logger the prepend every log entry with keyvals. */
func (l *TClientLogger) With(keyvals ...interface{}) log.Logger {
	logger := &TClientLogger{
		logger: l.logger,
	}

	if l.globalKeyvals != "" {
		logger.globalKeyvals = l.globalKeyvals + " "
	}

	logger.globalKeyvals += strings.TrimSuffix(fmt.Sprintln(keyvals...), "\n")

	return logger
}

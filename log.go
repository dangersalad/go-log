// Package log provides all the logging you need.
//
// See https://dave.cheney.net/2015/11/05/lets-talk-about-logging for
// rationale.
//
// Debug logging is controlled via environment variables. Set
// DEPLOY_ENV to "dev" or set LOG_DEBUG to a non empty value to enable
// the debug log.
package log // import "github.com/dangersalad/go-log"

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	debugPrefix   = "DBG"
	infoPrefix    = "NFO"
	debugEnabled  = false
	defaultLogger = NewLogger("main")
)

func init() {
	debugEnabled = checkDebugEnabled()
}

// Debug logs a debug message
func Debug(a ...interface{}) {
	defaultLogger.Debug(a...)
}

// Debugf logs a formatted debug message
func Debugf(f string, a ...interface{}) {
	defaultLogger.Debugf(f, a...)
}

// Info logs a message
func Info(a ...interface{}) {
	defaultLogger.Info(a...)
}

// Infof logs a formatted message
func Infof(f string, a ...interface{}) {
	defaultLogger.Infof(f, a...)
}

// Print is an alias for Info
func Print(a ...interface{}) {
	Info(a...)
}

// Printf is an alias for Infof
func Printf(f string, a ...interface{}) {
	Infof(f, a...)
}

// Logger is a logger with a prefix
type Logger struct {
	prefix string
}

// NewLogger returns a logger with the specified prefix
func NewLogger(p string) *Logger {
	return &Logger{
		prefix: p,
	}
}

// Debug logs a debug message with the logger's prefix
func (l *Logger) Debug(a ...interface{}) {
	debug(l.prefix, a...)
}

// Debugf logs a formatted debug message with the logger's prefix
func (l *Logger) Debugf(f string, a ...interface{}) {
	debugf(l.prefix, f, a...)
}

// Info logs a message with the logger's prefix
func (l *Logger) Info(a ...interface{}) {
	info(l.prefix, a...)
}

// Infof logs a formatted info message with the logger's prefix
func (l *Logger) Infof(f string, a ...interface{}) {
	infof(l.prefix, f, a...)
}

// Print is an alias for Info
func (l *Logger) Print(a ...interface{}) {
	l.Info(a...)
}

// Printf is an alias for Infof
func (l *Logger) Printf(f string, a ...interface{}) {
	l.Infof(f, a...)
}

func debug(prefix string, a ...interface{}) {
	if !debugEnabled {
		return
	}
	output(debugPrefix, prefix, a...)
}

func debugf(prefix, f string, a ...interface{}) {
	if !debugEnabled {
		return
	}
	outputf(debugPrefix, prefix, f, a...)
}

func info(prefix string, a ...interface{}) {
	output(infoPrefix, prefix, a...)
}

func infof(prefix, f string, a ...interface{}) {
	outputf(infoPrefix, prefix, f, a...)
}

func output(levelPrefix, loggerPrefix string, a ...interface{}) {
	if debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%s  | ", getCaller())}, a...)
	}
	a = append([]interface{}{fmt.Sprintf("%s  | ", loggerPrefix)}, a...)
	if debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%s  | ", levelPrefix)}, a...)
	}
	a = append([]interface{}{fmt.Sprintf("%s  | ", getTimestamp())}, a...)

	fmt.Println(a...)
}

func outputf(levelPrefix, loggerPrefix, f string, a ...interface{}) {
	if debugEnabled {
		f = fmt.Sprintf("%s  |  %s  |  %s  |  %s  |  %s", getTimestamp(), levelPrefix, loggerPrefix, getCaller(), f)
	} else {
		f = fmt.Sprintf("%s  |  %s  |  %s", getTimestamp(), loggerPrefix, f)
	}

	if f[len(f)-1] != '\n' {
		f += "\n"
	}
	fmt.Printf(f, a...)
}

func getTimestamp() string {
	return fmt.Sprintf("%-30s", time.Now().UTC().Format(time.RFC3339Nano))
}

func getCaller(s ...int) string {
	skip := 0
	if len(s) == 1 {
		skip = s[0]
	}
	_, fullfile, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	if strings.Index(fullfile, "log.go") >= 0 || strings.Index(fullfile, "logging.go") >= 0 {
		return getCaller(skip + 2)
	}

	parts := strings.Split(fullfile, "/")
	file := strings.Join(parts[len(parts)-2:], "/")

	return fmt.Sprintf("%s:%d", file, line)
}

func checkDebugEnabled() bool {
	return os.Getenv("DEPLOY_ENV") == "dev" || os.Getenv("LOG_DEBUG") != ""
}
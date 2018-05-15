// Package log provides all the logging you need.
//
// See https://dave.cheney.net/2015/11/05/lets-talk-about-logging for
// rationale.
//
// Debug logging is controlled via environment variables. Set
// DEPLOY_ENV to "dev" or "development", or set LOG_DEBUG to a non
// empty value to enable the debug log.
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
	defaultLogger = NewLogger("main", true)
)

// Debug logs a debug message
func Debug(a ...interface{}) {
	defaultLogger.Debug(a...)
}

// Debugln logs a debug message
func Debugln(a ...interface{}) {
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

// Infoln logs a message
func Infoln(a ...interface{}) {
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

// Println is an alias for Info
func Println(a ...interface{}) {
	Info(a...)
}

// Printf is an alias for Infof
func Printf(f string, a ...interface{}) {
	Infof(f, a...)
}

// Die will log out an error with "%+v" and exit the process with an
// optional code.
//
// Only available at the package level. This is meant to be used with
// github.com/pkg/errors and only at the top level of a process to
// handle errors that bubble up.
func Die(err error, code ...int) {
	defaultLogger.die(err, code...)
}

// Logger is a logger with a prefix
type Logger struct {
	prefix       string
	debugEnabled bool
}

// NewLogger returns a logger with the specified prefix and debugging
// possibly enabled. If debugEnabled is `false`, debug logging is
// always disabled. If `true`, it will follow the environment
// variables.
func NewLogger(prefix string, debugEnabled bool) *Logger {
	d := false
	if debugEnabled {
		d = checkDebugEnabled()
	}
	return &Logger{
		prefix:       prefix,
		debugEnabled: d,
	}
}

// Debug logs a debug message with the logger's prefix
func (l *Logger) Debug(a ...interface{}) {
	l.debug(l.prefix, a...)
}

// Debugln logs a debug message with the logger's prefix
func (l *Logger) Debugln(a ...interface{}) {
	l.debug(l.prefix, a...)
}

// Debugf logs a formatted debug message with the logger's prefix
func (l *Logger) Debugf(f string, a ...interface{}) {
	l.debugf(l.prefix, f, a...)
}

// Info logs a message with the logger's prefix
func (l *Logger) Info(a ...interface{}) {
	l.info(l.prefix, a...)
}

// Infoln logs a message with the logger's prefix
func (l *Logger) Infoln(a ...interface{}) {
	l.info(l.prefix, a...)
}

// Infof logs a formatted info message with the logger's prefix
func (l *Logger) Infof(f string, a ...interface{}) {
	l.infof(l.prefix, f, a...)
}

// Print is an alias for Info
func (l *Logger) Print(a ...interface{}) {
	l.Info(a...)
}

// Println is an alias for Info
func (l *Logger) Println(a ...interface{}) {
	l.Info(a...)
}

// Printf is an alias for Infof
func (l *Logger) Printf(f string, a ...interface{}) {
	l.Infof(f, a...)
}

func (l *Logger) debug(prefix string, a ...interface{}) {
	if !l.debugEnabled {
		return
	}
	l.output(debugPrefix, prefix, a...)
}

func (l *Logger) debugf(prefix, f string, a ...interface{}) {
	if !l.debugEnabled {
		return
	}
	l.outputf(debugPrefix, prefix, f, a...)
}

func (l *Logger) info(prefix string, a ...interface{}) {
	l.output(infoPrefix, prefix, a...)
}

func (l *Logger) infof(prefix, f string, a ...interface{}) {
	l.outputf(infoPrefix, prefix, f, a...)
}

func (l *Logger) output(levelPrefix, loggerPrefix string, a ...interface{}) {
	if l.debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%s  | ", getCaller())}, a...)
	}
	a = append([]interface{}{fmt.Sprintf("%s  | ", loggerPrefix)}, a...)
	if l.debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%s  | ", levelPrefix)}, a...)
	}
	a = append([]interface{}{fmt.Sprintf("%s  | ", getTimestamp())}, a...)

	fmt.Println(a...)
}

func (l *Logger) outputf(levelPrefix, loggerPrefix, f string, a ...interface{}) {
	if l.debugEnabled {
		f = fmt.Sprintf("%s  |  %s  |  %s  |  %s  |  %s", getTimestamp(), levelPrefix, loggerPrefix, getCaller(), f)
	} else {
		f = fmt.Sprintf("%s  |  %s  |  %s", getTimestamp(), loggerPrefix, f)
	}

	if f[len(f)-1] != '\n' {
		f += "\n"
	}
	fmt.Printf(f, a...)
}

func (l *Logger) die(err error, code ...int) {
	fmt.Fprintf(os.Stderr, "DIE %s\n%+v\n", getTimestamp(), err)
	c := 1
	if len(code) > 0 {
		c = code[0]
	}
	os.Exit(c)
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
	deployEnv := os.Getenv("DEPLOY_ENV")
	return deployEnv == "dev" ||
		deployEnv == "development" ||
		os.Getenv("LOG_DEBUG") != ""
}

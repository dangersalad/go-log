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
)

var (
	debugPrefix   = "DBG"
	infoPrefix    = "NFO"
	defaultLogger = NewLogger("main", true)
)

// SetDefaultName changes the name of the package level logger.
func SetDefaultName(n string) {
	if len(n) > prefixLimit {
		defaultLogger.prefix = n[0:prefixLimit]
	} else {
		defaultLogger.prefix = n
	}
}

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

const prefixLimit = 6
const callerLimit = 22

// NewLogger returns a logger with the specified prefix and debugging
// possibly enabled. If debugEnabled is `false`, debug logging is
// always disabled. If `true`, it will follow the environment
// variables.
func NewLogger(prefix string, debugEnabled bool) *Logger {
	d := false
	if debugEnabled {
		d = checkDebugEnabled()
	}
	// limit prefix
	if len(prefix) > prefixLimit {
		prefix = prefix[0:prefixLimit]
	}
	return &Logger{
		prefix:       prefix,
		debugEnabled: d,
	}
}

// Debug logs a debug message with the logger's prefix
func (l *Logger) Debug(a ...interface{}) {
	l.debug(a...)
}

// Debugln logs a debug message with the logger's prefix
func (l *Logger) Debugln(a ...interface{}) {
	l.debug(a...)
}

// Debugf logs a formatted debug message with the logger's prefix
func (l *Logger) Debugf(f string, a ...interface{}) {
	l.debugf(f, a...)
}

// Info logs a message with the logger's prefix
func (l *Logger) Info(a ...interface{}) {
	l.info(a...)
}

// Infoln logs a message with the logger's prefix
func (l *Logger) Infoln(a ...interface{}) {
	l.info(a...)
}

// Infof logs a formatted info message with the logger's prefix
func (l *Logger) Infof(f string, a ...interface{}) {
	l.infof(f, a...)
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

func (l *Logger) debug(a ...interface{}) {
	if !l.debugEnabled {
		return
	}
	l.output(debugPrefix, a...)
}

func (l *Logger) debugf(f string, a ...interface{}) {
	if !l.debugEnabled {
		return
	}
	l.outputf(debugPrefix, f, a...)
}

func (l *Logger) info(a ...interface{}) {
	l.output(infoPrefix, a...)
}

func (l *Logger) infof(f string, a ...interface{}) {
	l.outputf(infoPrefix, f, a...)
}

func (l *Logger) output(levelPrefix string, a ...interface{}) {
	if l.debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%-22s  | ", getCaller())}, a...)
	}
	a = append([]interface{}{fmt.Sprintf("%-6s  | ", l.prefix)}, a...)
	if l.debugEnabled {
		a = append([]interface{}{fmt.Sprintf("%s  | ", levelPrefix)}, a...)
	}

	fmt.Println(a...)
}

func (l *Logger) outputf(levelPrefix, f string, a ...interface{}) {
	if l.debugEnabled {
		f = fmt.Sprintf("%s  |  %-6s  |  %-22s  |  %s", levelPrefix, l.prefix, getCaller(), f)
	} else {
		f = fmt.Sprintf("%-6s  |  %s", l.prefix, f)
	}

	if f[len(f)-1] != '\n' {
		f += "\n"
	}
	fmt.Printf(f, a...)
}

func (l *Logger) die(err error, code ...int) {
	fmt.Fprintf(os.Stderr, "DIE\n%+v\n", err)
	c := 1
	if len(code) > 0 {
		c = code[0]
	}
	os.Exit(c)
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

	return normalizeCaller(line, fullfile)

}

func normalizeCaller(line int, fullfile string, counts ...int) string {
	count := 4
	if len(counts) > 0 {
		count = counts[0]
	}
	parts := strings.Split(fullfile, "/")
	file := strings.Join(parts[len(parts)-count:], "/")

	caller := fmt.Sprintf("%s:%d", file, line)
	if len(caller) > callerLimit {
		if count == 1 {
			return caller[len(caller)-callerLimit:]
		}
		return normalizeCaller(line, fullfile, count-1)
	}

	return caller
}

func checkDebugEnabled() bool {
	deployEnv := os.Getenv("DEPLOY_ENV")
	return deployEnv == "dev" ||
		deployEnv == "development" ||
		deployEnv == "test" ||
		deployEnv == "testing" ||
		os.Getenv("LOG_DEBUG") != ""
}

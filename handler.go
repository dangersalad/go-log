package log

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"
)

// DefaultPathLogBlacklist is a basic set of paths to ignore for logging
var DefaultPathLogBlacklist = regexp.MustCompile(`/ping|/healthz`)

type statusWriter struct {
	http.ResponseWriter
	status int
	body   string
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	w.body = string(b)
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacking not supported")
}

// HTTPHandler returns a handler that will log out request data.
//
// If the logger is nil, the default "main" logger is
// used.
//
// blacklist can be nil, in which case all calls are logged
func HTTPHandler(h http.Handler, logger *Logger, blacklist *regexp.Regexp) http.Handler {

	if logger == nil {
		logger = defaultLogger
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{
			ResponseWriter: w,
			status:         200,
		}
		h.ServeHTTP(sw, r)
		// get the diff and parse that time
		diff := time.Now().Sub(start)
		// don't log for certain paths
		if blacklist != nil && blacklist.MatchString(r.URL.Path) {
			return
		}
		diffStr := diff.String()
		if diff > time.Second {
			diffStr = diff.Truncate(time.Millisecond).String()
		} else if diff > time.Millisecond {
			diffStr = fmt.Sprintf("%0.3fms", float64(diff.Nanoseconds())/10000000.0)
		}
		switch c := sw.status; true {
		case c >= 500:
			logger.Infof("%s %s [%d] (%s)", r.Method, r.URL, c, diffStr)
		default:
			logger.Debugf("%s %s [%d] (%s)", r.Method, r.URL, c, diffStr)
		}
	})
}

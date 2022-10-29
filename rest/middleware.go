package rest

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewMiddleware(log *zap.SugaredLogger, proxyForwardedHeader string, logAllRequests bool) *Middleware {
	sugaredLogger := log.Desugar()
	zapLog := sugaredLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	middleware := &Middleware{
		log:                  zapLog,
		proxyForwardedHeader: proxyForwardedHeader,
		logAllRequests:       logAllRequests,
	}
	return middleware
}

type Middleware struct {
	log                  *zap.SugaredLogger
	logAllRequests       bool
	proxyForwardedHeader string
}

func (m *Middleware) sessionMiddleware(next appHandlerFunc) appHandlerFunc {

	return func(ctx *RequestContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
		ctx.TraceID = uuid.New().String()
		context.Set(r, "ctx", *ctx)
		return next(ctx, w, r)
	}
}

func (m *Middleware) notProtectedMiddleware(next appHandlerFunc) appHandlerFunc {
	return func(ctx *RequestContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return next(ctx, w, r)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	hijacked   bool
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	if lrw.hijacked {
		return
	}
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("Application does not support Hijack")
	}
	lrw.hijacked = true
	return hijacker.Hijack()
}

func (m *Middleware) loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		lw := &loggingResponseWriter{
			w,
			false,
			http.StatusOK,
		}

		next.ServeHTTP(lw, r)

		var traceID string

		ctxI := context.Get(r, "ctx")
		if ctxI != nil {
			if ctx, ok := ctxI.(RequestContext); ok {
				traceID = ctx.TraceID
			}
		}

		errI := context.Get(r, "error")
		if errI != nil {
			if err, ok := errI.(error); ok {
				m.log.Infow("REQUEST", "trace_id", traceID, "error", err.Error())
			} else {
				msg := fmt.Sprintf("REQUEST: %v: Unknown error: %#v", traceID, errI)
				m.log.Error(msg)
			}
		} else {
			// skip traceID if no error
			traceID = ""
		}

		if m.logAllRequests || r.Method != http.MethodGet || errI != nil {
			m.log.Infof("REQUEST [%v] %s %v %v %v", time.Since(startTime), traceID, lw.statusCode, r.Method, r.RequestURI)
		}
	})
}

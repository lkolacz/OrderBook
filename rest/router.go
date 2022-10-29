package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lkolacz/OrderBook/rest/config"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

const (
	PathPrefix = "/api/v1"

	PathHealthcheckVersion = "/healthcheck/version"
	PathOrderHandling      = "/order-handling"
)

func WrapPathPrefix(uri string) string {
	return strings.Join([]string{PathPrefix, uri}, "")
}

type appHandlerFunc func(ctx *RequestContext, w http.ResponseWriter, r *http.Request) (interface{}, error)

func (a appHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := &RequestContext{}

	resp, err := a(ctx, w, r)

	if strings.ToLower(r.Header.Get("connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("upgrade")) == "websocket" {
		if err != nil {
			var apiErr *APIError

			if !errors.As(err, &apiErr) {
				apiErr = ErrorInternal().Wrap(err)
			}

			context.Set(r, "error", apiErr)
		}
		return
	}

	FormatJSONResp(w, r, resp, err)
}

type route struct {
	path    string
	method  string
	handler appHandlerFunc
}

type Router struct {
	log        *zap.SugaredLogger
	config     *config.Config
	router     http.Handler
	handler    *handler
	middleware *Middleware
	version    *Version
}

func NewQRouter(log *zap.SugaredLogger, config *config.Config, version *Version) *Router {

	rt := &Router{
		log:        log,
		config:     config,
		router:     nil,
		handler:    &handler{},
		middleware: NewMiddleware(log, config.HTTP.ProxyForwardedHeader, config.HTTP.LogAllRequests),
		version:    version,
	}

	rt.handler = &handler{
		cfg:     *config,
		log:     log,
		version: version,
	}
	rt.router = rt.SetHandlers()

	return rt
}

// set all handler (endpoints & websockets)
func (r *Router) SetHandlers() http.Handler {

	routes := []route{
		{PathHealthcheckVersion, http.MethodGet, r.handler.HealthCheckVersion},
		{PathOrderHandling, MethodWebsocket, r.handler.OrderHandling},
	}

	router := mux.NewRouter().PathPrefix(PathPrefix).Subrouter()

	for _, route := range routes {

		middle := r.middleware.notProtectedMiddleware

		if route.method == MethodWebsocket {
			router.Handle(route.path, r.middleware.sessionMiddleware(middle(route.handler)))
		} else {
			router.Handle(route.path, r.middleware.sessionMiddleware(middle(route.handler))).Methods(route.method)
		}
	}

	router.Use(r.middleware.loggingMiddleware)

	r.printRoutes(router)

	return r.setupCORS(router)
}

// Start Run the service
func (r *Router) Start() error {
	errChan := make(chan error)
	r.StartHTTPListener(errChan)

	return <-errChan
}

// StartHTTPListener run the HTTP listener
func (r *Router) StartHTTPListener(errChan chan error) {
	r.log.Infof("CORS policies: %s", strings.Join(r.config.HTTP.CORSAllowOrigins, ","))
	if r.config.HTTP.ProxyForwardedHeader != "" {
		r.log.Info("Use Proxy forwarded-for header: %s", r.config.HTTP.ProxyForwardedHeader)
	}
	r.log.Infof("Starting listener on %v", r.config.HTTP.Addr)

	errChan <- http.ListenAndServe(r.config.HTTP.Addr, context.ClearHandler(r.router))
}

func (r *Router) setupCORS(h http.Handler) http.Handler {
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
		handlers.AllowedOrigins(r.config.HTTP.CORSAllowOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD"}),
		handlers.AllowCredentials(),
	)
	return cors(h)
}

func (r *Router) printRoutes(router *mux.Router) {
	if err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		if templ, err := route.GetPathTemplate(); err == nil {
			if method, err := route.GetMethods(); err == nil {
				for _, m := range method {
					r.log.Debugf("Registered handler %v %v", m, templ)
				}
			} else {
				r.log.Debugf("Registered websocket %v", templ)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

// WriteErrorOnHttp writes the error response as JSON
func WriteErrorOnHttp(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *APIError

	if !errors.As(err, &apiErr) {
		apiErr = ErrorInternal().Wrap(err)
	}
	context.Set(r, "error", apiErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Code())
	w.Write(apiErr.JSON())
}

// FormatJSONResp encodes response as JSON and handle errors
func FormatJSONResp(w http.ResponseWriter, r *http.Request, v interface{}, err error) {
	if err != nil {
		WriteErrorOnHttp(w, r, err)
		return
	}

	if v == nil {
		v = &struct {
			Code int
			Msg  string
		}{
			Code: http.StatusOK,
			Msg:  http.StatusText(http.StatusOK),
		}
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		WriteErrorOnHttp(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

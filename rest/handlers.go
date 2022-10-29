package rest

import (
	"net/http"

	"github.com/lkolacz/OrderBook/rest/config"
	"github.com/lkolacz/OrderBook/rest/core"

	"go.uber.org/zap"
)

type handler struct {
	core    core.AppCore
	cfg     config.Config
	log     *zap.SugaredLogger
	version *Version
}

// HealthCheckVersion
//
// swagger:route GET /healthcheck/version
//
// Check application version.
//
func (h *handler) HealthCheckVersion(_ *RequestContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	h.log.Debug("Handler for HealthCheckVersion endpoint")
	w.Header().Set("Content-Type", "application/json")
	response := h.version
	return response, nil
}

// OrderHandling
//
// swagger:route POST /order-handling  orderHandling OrderHandling
//
// Get or Set transation on Order Book via WebSocket.
//
func (h *handler) OrderHandling(_ *RequestContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	h.log.Debug("Handler for OrderHandling endpoint")
	WebSocketOrderHandling(h, w, r)
	return nil, nil
}

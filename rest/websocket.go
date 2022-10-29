package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lkolacz/OrderBook/rest/api"
	"github.com/lkolacz/OrderBook/rest/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var VisibleStock = 2

func WebSocketOrderHandling(h *handler, w http.ResponseWriter, r *http.Request) {
	h.log.Debug("Handler for WebSocketOrderHandling")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// current status
	session := core.GetSessionOrders(VisibleStock)
	err = ws.WriteJSON(session)
	if err != nil {
		h.log.Errorf("error when writting to client: ", err)
	}

	for {
		var v api.Parser = &api.Errand{}

		if err := ws.ReadJSON(v); err != nil {
			h.log.Errorf("error when reading from client: ", err)
			continue
		}
		h.log.Debugf("read message from client: %v", v.Parse())

		transations, err := core.ProcessErrand(v.Parse())
		if err != nil {
			h.log.Errorf("error when process errand to client: ", err)
		}
		session := core.GetSessionOrders(VisibleStock)

		err = ws.WriteJSON(session)
		if err != nil {
			h.log.Errorf("error when writting to client: ", err)
		}
		for _, tran := range transations {
			err = ws.WriteJSON(tran)
			if err != nil {
				h.log.Errorf("error when writting to client: ", err)
			}
		}
	}

}

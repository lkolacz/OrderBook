package api

import "encoding/json"

const (
	ERANDTYPELIMIT   = "Limit"
	ERANDTYPEICEBERG = "Iceberg"
	DIRECTIONBUY     = "Buy"
	DIRECTIONSELL    = "Sell"
)

type OrderMsg struct {
	Direction string `json:"direction"`
	ID        int    `json:"id"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
	Peak      int    `json:"peak,omitempty"`
}

type OrderBase struct {
	ID       int `json:"id"`
	Price    int `json:"price"`
	Quantity int `json:"quantity"`
}

type Errand struct {
	Type  string   `json:"type"`
	Order OrderMsg `json:"order"`
}

type Parser interface {
	Parse() string
}

func (a Errand) Parse() string {
	out, _ := json.Marshal(a)
	return string(out)
}

type Transation struct {
	BuyOrderId  int `json:"buyOrderId"`
	SellOrderId int `json:"sellOrderId"`
	Price       int `json:"price"`
	Quantity    int `json:"quantity"`
}

type SessionOrders struct {
	BuyOrders  []OrderBase `json:"buyOrders"`
	SellOrders []OrderBase `json:"sellOrders"`
}

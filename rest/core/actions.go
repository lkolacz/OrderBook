package core

import (
	"encoding/json"
	"math"
	"sort"
	"sync"

	"github.com/lkolacz/OrderBook/rest/api"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	mutex         = sync.Mutex{}
	orderBookSell = []api.OrderMsg{}
	orderBookBuy  = []api.OrderMsg{}
)

type ByPrice map[int]api.OrderMsg

func (a ByPrice) Len() int           { return len(a) }
func (a ByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }

func insertByPrice(arr []api.OrderMsg, elem api.OrderMsg, descending ...bool) []api.OrderMsg {
	var (
		arrLen, i int
		desc      bool
	)
	arrLen = len(arr)
	if len(descending) > 0 {
		desc = descending[0]
	} else {
		desc = false
	}
	if desc {
		i = sort.Search(arrLen, func(i int) bool { return arr[i].Price < elem.Price })
	} else {
		i = sort.Search(arrLen, func(i int) bool { return arr[i].Price > elem.Price })
	}

	if i < arrLen {
		// extend arr to new size (add new empty item)
		arr = append(arr, api.OrderMsg{})
		// Copy over elements sourced from index i, into elements starting at index i+1.
		copy(arr[i+1:], arr[i:])
		arr[i] = elem
	} else {
		arr = append(arr, elem)
	}
	return arr
}

func ProcessErrand(income string) ([]api.Transation, error) {

	var (
		orderBookBuyLen  int
		orderBookSellLen int
		transations      []api.Transation
		errand           = &api.Errand{}
	)
	err := json.Unmarshal([]byte(income), errand)
	if err != nil {
		return transations, err
	}

	caser := cases.Title(language.Und)
	direction := caser.String(errand.Order.Direction)

	// Lock from now the resource and unlock it to the end of function.
	mutex.Lock()
	defer mutex.Unlock()

	if direction == api.DIRECTIONBUY {
		// buy
		orderBookBuy = insertByPrice(orderBookBuy, errand.Order, true)
		orderBookSellLen = len(orderBookSell)

	loopBuyFromSell:
		for i := 0; i < orderBookSellLen && orderBookSellLen > 0; i++ {
			breakAtTheEnd := false
			if orderBookSell[i].Price <= orderBookBuy[0].Price {
				sellQuan := getRealQuantity(orderBookSell[i])
				buyQuan := getRealQuantity(orderBookBuy[0])
				realQuanOfSell := int(math.Min(float64(sellQuan), float64(buyQuan)))

				orderBookSell[i].Quantity -= realQuanOfSell
				orderBookBuy[0].Quantity -= realQuanOfSell

				transations = appendTransations(transations, realQuanOfSell, orderBookBuy[0], orderBookSell[i])

				didDelete := false
				if orderBookBuy[0].Quantity <= 0 {
					// delete this index - its empty
					orderBookBuy = slices.Delete(orderBookBuy, 0, 1)
					// nothing more to sell, goodbey
					breakAtTheEnd = true
				}
				if orderBookSell[i].Quantity <= 0 {
					// delete this index - its empty
					orderBookSell = slices.Delete(orderBookSell, 0, 1)
					orderBookSellLen -= 1
					didDelete = true
					// we don't want to break here - maybe next one will have a deal
				}

				if breakAtTheEnd {
					// no more operation at stock
					break loopBuyFromSell
				} else if didDelete {
					goto loopBuyFromSell
				} else if i == orderBookSellLen-1 {
					goto loopBuyFromSell
				} else if i+1 <= orderBookSellLen {
					if orderBookSell[i].Peak > 0 && orderBookSell[i].Price != orderBookSell[i+1].Price {
						goto loopBuyFromSell
					}
				}

			}
		} // end for
	} else {
		// sell
		orderBookSell = insertByPrice(orderBookSell, errand.Order)
		orderBookBuyLen = len(orderBookBuy)

	loopBuyFromBuy:
		for i := 0; i < orderBookBuyLen && orderBookBuyLen > 0; i++ {
			breakAtTheEnd := false
			if orderBookSell[0].Price <= orderBookBuy[i].Price {
				sellQuan := getRealQuantity(orderBookSell[0])
				buyQuan := getRealQuantity(orderBookBuy[i])
				realQuan := int(math.Min(float64(sellQuan), float64(buyQuan)))

				orderBookSell[0].Quantity -= realQuan
				orderBookBuy[i].Quantity -= realQuan

				transations = appendTransations(transations, realQuan, orderBookBuy[i], orderBookSell[0])

				didDelete := false
				if orderBookBuy[i].Quantity <= 0 {
					// delete this index - its empty
					orderBookBuy = slices.Delete(orderBookBuy, 0, 1)
					orderBookBuyLen -= 1
					didDelete = true
					// we don't want to break here - maybe next one will have a deal
				}
				if orderBookSell[0].Quantity <= 0 {
					// delete this index - its empty
					orderBookSell = slices.Delete(orderBookSell, 0, 1)
					// nothing more to sell, goodbey
					breakAtTheEnd = true
				}

				if breakAtTheEnd {
					// no more operation at stock
					break loopBuyFromBuy
				} else if didDelete {
					goto loopBuyFromBuy
				} else if i == orderBookBuyLen-1 {
					goto loopBuyFromBuy
				} else if i+1 <= orderBookBuyLen {
					if orderBookBuy[i].Peak > 0 && orderBookBuy[i].Price != orderBookBuy[i+1].Price {
						goto loopBuyFromBuy
					}
				}

			}
		} // end for

	}

	return transations, nil
}

func getOrders(arr []api.OrderMsg) []api.OrderBase {
	items := []api.OrderBase{}
	for _, elem := range arr {
		var icebergQuantity int
		if elem.Peak > 0 {
			if elem.Quantity > elem.Peak {
				icebergQuantity = elem.Peak
			} else {
				icebergQuantity = elem.Quantity
			}
		} else {
			icebergQuantity = elem.Quantity
		}
		items = append(items, api.OrderBase{ID: elem.ID, Price: elem.Price, Quantity: icebergQuantity})
	}
	return items
}

func GetSessionOrders(idx int) api.SessionOrders {
	var (
		BuyOrders  []api.OrderBase
		SellOrders []api.OrderBase
	)
	sizeBuy := len(orderBookBuy)
	if idx >= sizeBuy {
		BuyOrders = getOrders(orderBookBuy[:sizeBuy]) // already sorted ascending
	} else {
		BuyOrders = getOrders(orderBookBuy[:idx]) // already sorted ascending
	}

	sizeSell := len(orderBookSell)
	if idx >= sizeSell {
		SellOrders = getOrders(orderBookSell[:sizeSell]) // already sorted descending
	} else {
		SellOrders = getOrders(orderBookSell[:idx]) // already sorted descending
	}

	return api.SessionOrders{
		BuyOrders:  BuyOrders,
		SellOrders: SellOrders,
	}
}

func getRealQuantity(item api.OrderMsg) int {
	if item.Peak == 0 {
		return item.Quantity
	} else if item.Peak > item.Quantity {
		return item.Quantity
	}
	return item.Peak
}

func appendTransations(transations []api.Transation, quantity int, buy, sell api.OrderMsg) []api.Transation {
	tran := api.Transation{
		BuyOrderId:  buy.ID,
		SellOrderId: sell.ID,
		Price:       buy.Price,
		Quantity:    quantity,
	}
	return append(transations, tran)
}

func ClearStock() {
	orderBookSell = []api.OrderMsg{}
	orderBookBuy = []api.OrderMsg{}
}

package core_test

import (
	"testing"

	"github.com/lkolacz/OrderBook/rest/api"
	"github.com/lkolacz/OrderBook/rest/core"
	"github.com/stretchr/testify/assert"
)

func TestCoreAction(t *testing.T) {
	core.ClearStock()
	t.Run(
		"Core action - Process Errand (action on buy) no. 1",
		func(t *testing.T) {
			value1 := `{"type":"Limit","order":{"direction":"Buy","id":1,"price":14,"quantity":20}}`
			value2 := `{"type":"Iceberg","order":{"direction":"Buy","id":2,"price":15,"quantity":60,"peak":20}}`
			value3 := `{"type":"Limit","order":{"direction":"Sell","id":3,"price":16,"quantity":15}}`
			value4 := `{"type":"Limit","order":{"direction":"Buy","id":4,"price":16,"quantity":15,"peak":10}}`

			trans, err := core.ProcessErrand(value1)
			assert.NoError(t, err)
			assert.Empty(t, trans)
			trans, err = core.ProcessErrand(value2)
			assert.NoError(t, err)
			assert.Empty(t, trans)
			trans, err = core.ProcessErrand(value3)
			assert.NoError(t, err)
			assert.Empty(t, trans)
			trans, err = core.ProcessErrand(value4)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans := []api.Transation{
				{
					BuyOrderId:  4,
					SellOrderId: 3,
					Price:       16,
					Quantity:    10,
				},
				{
					BuyOrderId:  4,
					SellOrderId: 3,
					Price:       16,
					Quantity:    5,
				}}

			assert.EqualValues(t, expectedTrans, trans)
			session := core.GetSessionOrders(2)
			assert.Empty(t, session.SellOrders)
			assert.Equal(t, 2, len(session.BuyOrders))

		})

	t.Run(
		"Core action - Process Errand (action on buy) no. 2",
		func(t *testing.T) {
			// current state:
			// "buyOrders": [
			// 	{
			// 		"id": 2,
			// 		"price": 15,
			// 		"quantity": 20
			// 	},
			// 	{
			// 		"id": 1,
			// 		"price": 14,
			// 		"quantity": 20
			// 	}
			// ],
			// "sellOrders": [
			// ]
			//value1 := `{"type":"Limit","order":{"direction":"Buy","id":1,"price":14,"quantity":20}}`
			//value2 := `{"type":"Iceberg","order":{"direction":"Buy","id":2,"price":15,"quantity":60,"peak":20}}`
			value5 := `{"type":"Limit","order":{"direction":"Sell","id":5,"price":16,"quantity":10}}`
			value6 := `{"type":"Limit","order":{"direction":"Buy","id":6,"price":16,"quantity":15}}`

			trans, err := core.ProcessErrand(value5)
			assert.NoError(t, err)
			assert.Empty(t, trans)
			trans, err = core.ProcessErrand(value6)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans := []api.Transation{
				{
					BuyOrderId:  6,
					SellOrderId: 5,
					Price:       16,
					Quantity:    10,
				}}

			assert.EqualValues(t, expectedTrans, trans)
			session := core.GetSessionOrders(2)
			assert.Empty(t, session.SellOrders)
			assert.Equal(t, 2, len(session.BuyOrders))
			assert.Equal(t, 5, session.BuyOrders[0].Quantity)
			assert.Equal(t, 6, session.BuyOrders[0].ID)
			assert.Equal(t, 16, session.BuyOrders[0].Price)

		})

	t.Run(
		"Core action - Process Errand (action on sell) no. 3",
		func(t *testing.T) {
			//value1 := `{"type":"Limit","order":{"direction":"Buy","id":1,"price":14,"quantity":20}}`
			//value2 := `{"type":"Iceberg","order":{"direction":"Buy","id":2,"price":15,"quantity":60,"peak":20}}`
			//value6 := `{"type":"Limit","order":{"direction":"Buy","id":6,"price":16,"quantity":5}}`
			value7 := `{"type":"Limit","order":{"direction":"Sell","id":7,"price":15,"quantity":25}}`

			trans, err := core.ProcessErrand(value7)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans := []api.Transation{
				{
					BuyOrderId:  6,
					SellOrderId: 7,
					Price:       16,
					Quantity:    5,
				},
				{
					BuyOrderId:  2,
					SellOrderId: 7,
					Price:       15,
					Quantity:    20,
				}}

			assert.EqualValues(t, expectedTrans, trans)
			session := core.GetSessionOrders(2)
			assert.Empty(t, session.SellOrders)
			assert.Equal(t, 2, len(session.BuyOrders))
			assert.Equal(t, 20, session.BuyOrders[0].Quantity)
			assert.Equal(t, 2, session.BuyOrders[0].ID)
			assert.Equal(t, 15, session.BuyOrders[0].Price)
		})

	t.Run(
		"Core action - Process Errand (action on buy) no. 4",
		func(t *testing.T) {
			//value := `{"type":"Limit","order":{"direction":"Buy","id":1,"price":14,"quantity":20}}`
			//value1 := `{"type":"Iceberg","order":{"direction":"Buy","id":2,"price":15,"quantity":40,"peak":20}}`
			value11 := `{"type": "Iceberg", "order": {"direction": "Sell", "id": 11, "price": 100, "quantity": 200, "peak": 100}}`
			value12 := `{"type": "Iceberg", "order": {"direction": "Sell", "id": 12, "price": 100, "quantity": 300, "peak": 100}}`
			value13 := `{"type": "Iceberg", "order": {"direction": "Sell", "id": 13, "price": 100, "quantity": 200, "peak": 100}}`
			value14 := `{"type": "Iceberg", "order": {"direction": "Buy", "id": 14, "price": 100, "quantity": 500, "peak": 100}}`

			for _, val := range []string{value11, value12, value13} {
				trans, err := core.ProcessErrand(val)
				assert.NoError(t, err)
				assert.Empty(t, trans)
			}

			trans, err := core.ProcessErrand(value14)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)
			assert.Equal(t, 5, len(trans))
			sellOrdersId := 10
			for _, val := range trans {
				sellOrdersId += 1
				if sellOrdersId > 13 {
					sellOrdersId = 11
				}
				assert.Equal(t, sellOrdersId, val.SellOrderId)
				assert.Equal(t, 14, val.BuyOrderId)
				assert.Equal(t, 100, val.Quantity)
				assert.Equal(t, 100, val.Price)
			}

			session := core.GetSessionOrders(2)
			assert.Equal(t, 2, len(session.SellOrders))
			assert.Equal(t, 100, session.SellOrders[0].Quantity)
			assert.Equal(t, 12, session.SellOrders[0].ID)
			assert.Equal(t, 100, session.SellOrders[0].Price)
			assert.Equal(t, 100, session.SellOrders[1].Quantity)
			assert.Equal(t, 13, session.SellOrders[1].ID)
			assert.Equal(t, 100, session.SellOrders[1].Price)
		})
}

func TestCoreActionOtherScenario(t *testing.T) {
	core.ClearStock()
	t.Run(
		"Core action - Process Errand - insert buys, next insert sell",
		func(t *testing.T) {
			value1 := `{"type": "Limit", "order": {"direction": "Buy", "id": 1, "price": 14, "quantity": 20}}`
			value2 := `{"type": "Iceberg", "order": {"direction": "Buy", "id": 2, "price": 16, "quantity": 50, "peak": 20}}`
			value3 := `{"type": "Limit", "order": {"direction": "Buy", "id": 3, "price": 12, "quantity": 11}}`
			value4 := `{"type": "Limit", "order": {"direction": "Buy", "id": 4, "price": 13, "quantity": 5}}`
			value5 := `{"type": "Limit", "order": {"direction": "Buy", "id": 5, "price": 14, "quantity": 5}}`

			value6 := `{"type": "Limit", "order": {"direction": "Sell", "id": 6, "price": 20, "quantity": 200}}`
			value7 := `{"type": "Limit", "order": {"direction": "Sell", "id": 7, "price": 13, "quantity": 60}}`
			value8 := `{"type": "Limit", "order": {"direction": "Sell", "id": 8, "price": 13, "quantity": 10}}`
			value9 := `{"type": "Limit", "order": {"direction": "Sell", "id": 9, "price": 13, "quantity": 10}}`
			value10 := `{"type": "Limit", "order": {"direction": "Sell", "id": 10, "price": 12, "quantity": 15}}`

			for _, val := range []string{value1, value2, value3, value4, value5, value6} {
				trans, err := core.ProcessErrand(val)
				assert.NoError(t, err)
				assert.Empty(t, trans)
			}

			trans, err := core.ProcessErrand(value7)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans := []api.Transation{
				{
					BuyOrderId:  2,
					SellOrderId: 7,
					Price:       16,
					Quantity:    20,
				},
				{
					BuyOrderId:  2,
					SellOrderId: 7,
					Price:       16,
					Quantity:    20,
				},
				{
					BuyOrderId:  2,
					SellOrderId: 7,
					Price:       16,
					Quantity:    10,
				},
				{
					BuyOrderId:  1,
					SellOrderId: 7,
					Price:       14,
					Quantity:    10,
				},
			}
			assert.EqualValues(t, expectedTrans, trans)
			session := core.GetSessionOrders(2)
			assert.Equal(t, 2, len(session.BuyOrders))
			assert.Equal(t, 1, len(session.SellOrders))

			ids := []int{1, 5}
			price := []int{14, 14}
			quantity := []int{10, 5}
			for i, val := range session.BuyOrders {
				assert.Equal(t, ids[i], val.ID)
				assert.Equal(t, quantity[i], val.Quantity)
				assert.Equal(t, price[i], val.Price)
			}
			assert.Equal(t, 6, session.SellOrders[0].ID)
			assert.Equal(t, 200, session.SellOrders[0].Quantity)
			assert.Equal(t, 20, session.SellOrders[0].Price)

			trans, err = core.ProcessErrand(value8)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans = []api.Transation{
				{
					BuyOrderId:  1,
					SellOrderId: 8,
					Price:       14,
					Quantity:    10,
				}}
			assert.EqualValues(t, expectedTrans, trans)

			session = core.GetSessionOrders(2)
			assert.Equal(t, 2, len(session.BuyOrders))
			assert.Equal(t, 1, len(session.SellOrders))
			ids = []int{5, 4}
			price = []int{14, 13}
			quantity = []int{5, 5}
			for i, val := range session.BuyOrders {
				assert.Equal(t, ids[i], val.ID)
				assert.Equal(t, quantity[i], val.Quantity)
				assert.Equal(t, price[i], val.Price)
			}
			assert.Equal(t, 6, session.SellOrders[0].ID)
			assert.Equal(t, 200, session.SellOrders[0].Quantity)
			assert.Equal(t, 20, session.SellOrders[0].Price)

			trans, err = core.ProcessErrand(value9)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)

			expectedTrans = []api.Transation{
				{
					BuyOrderId:  5,
					SellOrderId: 9,
					Price:       14,
					Quantity:    5,
				},
				{
					BuyOrderId:  4,
					SellOrderId: 9,
					Price:       13,
					Quantity:    5,
				}}
			assert.EqualValues(t, expectedTrans, trans)

			session = core.GetSessionOrders(2)
			assert.Equal(t, 1, len(session.BuyOrders))
			assert.Equal(t, 1, len(session.SellOrders))

			assert.Equal(t, 3, session.BuyOrders[0].ID)
			assert.Equal(t, 11, session.BuyOrders[0].Quantity)
			assert.Equal(t, 12, session.BuyOrders[0].Price)

			assert.Equal(t, 6, session.SellOrders[0].ID)
			assert.Equal(t, 200, session.SellOrders[0].Quantity)
			assert.Equal(t, 20, session.SellOrders[0].Price)

			trans, err = core.ProcessErrand(value10)
			assert.NoError(t, err)
			assert.NotEmpty(t, trans)
			expectedTrans = []api.Transation{
				{
					BuyOrderId:  3,
					SellOrderId: 10,
					Price:       12,
					Quantity:    11,
				}}
			assert.EqualValues(t, expectedTrans, trans)

			session = core.GetSessionOrders(2)
			assert.Empty(t, session.BuyOrders)
			assert.Equal(t, 2, len(session.SellOrders))

			ids = []int{10, 6}
			price = []int{12, 20}
			quantity = []int{4, 200}
			for i, val := range session.SellOrders {
				assert.Equal(t, ids[i], val.ID)
				assert.Equal(t, quantity[i], val.Quantity)
				assert.Equal(t, price[i], val.Price)
			}

		})
}

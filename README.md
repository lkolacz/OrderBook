# OrderBook
Order Book (Limit & Iceberg Order)

# Install
On local machine download the code (git clone) and then run it locally:
```bash
% sh ./build.sh 
% ./out/order-book version
% go test ./tests/core/actions_test.go -v
% ./out/order-book start
```
or You can do it via docker image:
```bash
sh ./build.sh docker
% docker run -p 8007:8007 order-book:dev
```

# Use case
First, when You run the server, please see if it's working [here](http://127.0.0.1:8007/api/v1/healthcheck/version).

If Yes, then please use Postman or other tool that will help You connect via websocket to the endpoint: `127.0.0.1:8007/api/v1/order-handling`.

To see how it's working You can use this JSON data as follow (one by one, each line alone):

```JSON
{"type": "Limit", "order": {"direction": "Buy", "id": 1, "price": 14, "quantity": 20}}
{"type": "Iceberg", "order": {"direction": "Buy", "id": 2, "price": 16, "quantity": 50, "peak": 20}}
{"type": "Limit", "order": {"direction": "Buy", "id": 3, "price": 12, "quantity": 11}}
{"type": "Limit", "order": {"direction": "Buy", "id": 4, "price": 13, "quantity": 5}}
{"type": "Limit", "order": {"direction": "Buy", "id": 5, "price": 14, "quantity": 5}}
{"type": "Limit", "order": {"direction": "Sell", "id": 6, "price": 20, "quantity": 200}}
{"type": "Limit", "order": {"direction": "Sell", "id": 7, "price": 13, "quantity": 60}}
{"type": "Limit", "order": {"direction": "Sell", "id": 8, "price": 13, "quantity": 10}}
{"type": "Limit", "order": {"direction": "Sell", "id": 9, "price": 13, "quantity": 10}}
{"type": "Limit", "order": {"direction": "Sell", "id": 10, "price": 16, "quantity": 20}}
```

As a result You will get following output:

```JSON
{"buyOrders":[{"id":3,"price":12,"quantity":11}],"sellOrders":[{"id":10,"price":16,"quantity":20},{"id":6,"price":20,"quantity":200}]}
{"buyOrderId":4,"sellOrderId":9,"price":13,"quantity":5}
{"buyOrderId":5,"sellOrderId":9,"price":14,"quantity":5}
{"buyOrders":[{"id":3,"price":12,"quantity":11}],"sellOrders":[{"id":6,"price":20,"quantity":200}]}
{"buyOrderId":1,"sellOrderId":8,"price":14,"quantity":10}
{"buyOrders":[{"id":5,"price":14,"quantity":5},{"id":4,"price":13,"quantity":5}],"sellOrders":[{"id":6,"price":20,"quantity":200}]}
{"buyOrderId":1,"sellOrderId":7,"price":14,"quantity":10}
{"buyOrderId":2,"sellOrderId":7,"price":16,"quantity":10}
{"buyOrderId":2,"sellOrderId":7,"price":16,"quantity":20}
{"buyOrderId":2,"sellOrderId":7,"price":16,"quantity":20}
{"buyOrders":[{"id":1,"price":14,"quantity":10},{"id":5,"price":14,"quantity":5}],"sellOrders":[{"id":6,"price":20,"quantity":200}]}
{"buyOrders":[{"id":2,"price":16,"quantity":20},{"id":1,"price":14,"quantity":20}],"sellOrders":[{"id":6,"price":20,"quantity":200}]}
{"buyOrders":[{"id":2,"price":16,"quantity":20},{"id":1,"price":14,"quantity":20}],"sellOrders":[]}
{"buyOrders":[{"id":2,"price":16,"quantity":20},{"id":1,"price":14,"quantity":20}],"sellOrders":[]}
{"buyOrders":[{"id":2,"price":16,"quantity":20},{"id":1,"price":14,"quantity":20}],"sellOrders":[]}
{"buyOrders":[{"id":2,"price":16,"quantity":20},{"id":1,"price":14,"quantity":20}],"sellOrders":[]}
{"buyOrders":[{"id":1,"price":14,"quantity":20}],"sellOrders":[]}
{"buyOrders":[],"sellOrders":[]}
```

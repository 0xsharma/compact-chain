
# Compact-Chain
This is a simple implentation of Blockchain in Golang. This project is undertaken to enhance my knowlege on advanced(p2p, db, state-transitions) blockchain concepts. This can be used by anyone who wants to learn about blockchain and how it works. 

### Prerequisites

What things you need to install the software and how to install them.

```
Golang 1.20+
```

### Run Demo Chain

```
go mod tidy
go run main.go demo
```

### Run Multiple Nodes Chain

```
go mod tidy
go run main.go start <NODE_ID>

for example : 
terminal/instance 1 : go run main.go start 1
terminal/instance 2 : go run main.go start 2
terminal/instance 3 : go run main.go start 3
And so on....
```

### Send Transactions


Make sure a node is running and note the endpoint.

```
go run main.go send-tx --to <TO_ADDR> --privatekey <SENDER_PRIV_KEY> --value <TX_VALUE> --rpc <RPC_ADDR> --nonce <NONCE>
```
example (also, will run with fresh chain and default config) :
```
go run main.go send-tx --to 0x93a63fc45341fc02ac9cce62cc5aeb5c5799403e --privatekey c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a6 --value 1 --rpc localhost:17111 --nonce 0
```
(increase the nonce for the consecutive transactions by 1 to fire more transactions)
###### NOTE : Transactions can also be send using RPC calls directly.

### Run Tests

```
make test
```

### Modules Implemented

```
- Consensus (POW)
- p2p (gRPC)
- DbStore
- State Executor
- RPC (add and get Transactions)
- TxPool
- Encoding
- Hashing
```
### License
The entire code is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).


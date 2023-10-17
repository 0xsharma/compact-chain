
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


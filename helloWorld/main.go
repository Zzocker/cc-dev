package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// chaincode must implement two function
// 1. Init(stub shim.ChaincodeStubInterface)peer.Response
// 2. Invoke(stub shim.ChaincodeStubInterface)peer.Response
type chaincode struct{}

var (
	// mapping of supported methods
	// methodName => method implantation
	supportedMethod = map[string]func(stub shim.ChaincodeStubInterface, args []string) peer.Response{
		"hello": hello,
	}
)

func (chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// get methodName and argument sent by client
	// args by client : [ <funcName> , ...<string> ]
	methodName, args := stub.GetFunctionAndParameters()
	method, ok := supportedMethod[methodName]
	// check if method exists or not
	if !ok {
		return shim.Error(fmt.Sprintf(`method "<%s>" not supported`, methodName))
	}
	// invoke `methodName` method
	return method(stub, args)
}

func main() {
	if err := shim.Start(chaincode{}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func hello(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : [name]
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("require 1 [name] argument, but got %d", len(args)))
	}
	msg := fmt.Sprintf("Hello %s", args[0])
	return shim.Success([]byte(msg))
}

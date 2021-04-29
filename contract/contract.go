package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

var (
	// list of supported methods by invoke function
	// mapping of function name with actual function
	supportedMethods = map[string]func(stub shim.ChaincodeStubInterface, args []string) peer.Response{
		"create":   create,
		"query":    query,
		"transfer": transfer,
		"purge":    purge,
	}
)

// Chaincode : implements shim.Chaincode interface
type Chaincode struct{}

// Init : invoke method called with -I (--isInit) flag.
// use for initializing its initial state
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke : method called for updating and query state
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// get called method name along with arguments
	methodName, args := stub.GetFunctionAndParameters()

	// get method implantation for `methodName` from supported method mapping
	method, ok := supportedMethods[methodName]
	// check if method is supported or not
	if !ok {
		return shim.Error(fmt.Sprintf("method %s is not supported", methodName))
	}
	// call the method function
	return method(stub, args)
}

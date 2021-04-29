package contract

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestCreate(t *testing.T) {
	// create mock stub with given chaincode
	stub := shimtest.NewMockStub("_car_", new(Chaincode))

	// make mock invoke with function "create"
	// input : <Tx ID > , [][]byte{} - array of argument with 0th index as method name
	result := stub.MockInvoke("tx1", [][]byte{[]byte("create")})

	// result is a protocol buffer message with fields
	// - Message : string ,
	// - Status : int,
	// Payload : []byte
	if result.Status != shim.OK {
		t.Error("failed to invoke create a new car")
	}
}

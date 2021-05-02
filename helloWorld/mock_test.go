package main

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestHello(t *testing.T) {
	// create a mock stub
	stub := shimtest.NewMockStub("__hello__", chaincode{})

	name := "Pritam"

	result := stub.MockInvoke("txID", [][]byte{[]byte("hello"), []byte(name)})

	if result.Status != shim.OK {
		t.Errorf("error occurred : %s", result.Message)
	}

	t.Logf("Result : %s", string(result.Payload))
}

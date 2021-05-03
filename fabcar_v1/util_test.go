package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-protos-go/msp"
)

type mockCaller struct {
	pemFilePath string
	commonName  string
}

func TestGetCallerID(t *testing.T) {
	stub := shimtest.NewMockStub("__", nil)

	caller := mockCaller{
		pemFilePath: "./mockData/customer1_cert",
		commonName:  "customer1",
	}

	// set caller certifacts
	setCaller(stub, caller.readPEM())
	result := getCallerID(stub)

	t.Logf("caller's commonName returnred from cc : %s", result)
}

func setCaller(stub *shimtest.MockStub, certPEM string) {
	identity := &msp.SerializedIdentity{
		IdBytes: []byte(certPEM),
	}
	b, _ := proto.Marshal(identity)
	stub.Creator = b
}

func (m mockCaller) readPEM() string {
	file, _ := os.Open(m.pemFilePath)
	defer file.Close()
	pem, _ := ioutil.ReadAll(file)
	return string(pem)
}

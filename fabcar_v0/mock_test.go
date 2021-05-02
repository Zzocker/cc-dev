package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

var (
	mockCar = Car{
		Make:   "Toyota",
		Model:  "Prius",
		Colour: "Blue",
		Owner:  "Pritam",
	}

	mockCarID = "car1"
)

func TestCreate(t *testing.T) {
	// create a new mock stub
	stub := shimtest.NewMockStub("__car__", chaincode{})

	result := stub.MockInvoke("txID",
		[][]byte{
			[]byte("create"),  // function name
			[]byte(mockCarID), // args 0th index
			[]byte(mockCar.Make),
			[]byte(mockCar.Model),
			[]byte(mockCar.Colour),
			[]byte(mockCar.Owner),
		})

	if result.Status != shim.OK {
		t.Error("failed to create a new car")
	}

	// now check worldstate for carID
	carByte, err := stub.GetState(mockCarID)
	if err != nil {
		t.Error(err)
	}
	var car Car
	json.Unmarshal(carByte, &car)
	t.Logf("car from worldstate : %+v", car)
}

func TestTransfer(t *testing.T) {
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// insert mock car for transfering the car
	stub.MockTransactionStart("__tx")
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	newOwner := "newOwner"

	result := stub.MockInvoke("txID", [][]byte{
		[]byte("transfer"),
		[]byte(mockCarID),
		[]byte(newOwner),
	})
	if result.Status != shim.OK {
		t.Error("failed to transfer the car")
	}

	// now check worldstate for update carID
	carByte, err := stub.GetState(mockCarID)
	if err != nil {
		t.Error(err)
	}
	var car Car
	json.Unmarshal(carByte, &car)
	t.Logf("update car from worldstate : %+v", car)
}

func TestQuery(t *testing.T) {
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// insert mock car for querying the car
	stub.MockTransactionStart("__tx")
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	result := stub.MockInvoke("txID", [][]byte{
		[]byte("query"),
		[]byte(mockCarID),
	})
	if result.Status != shim.OK {
		t.Error("failed to query car")
	}
	var car Car
	json.Unmarshal(result.Payload, &car)
	t.Logf("car from chaincode : %+v", car)
}

func TestPurge(t *testing.T) {
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// insert mock car for deleting the car
	stub.MockTransactionStart("__tx")
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	result := stub.MockInvoke("txID",[][]byte{
		[]byte("purge"),
		[]byte(mockCarID),
	})

	if result.Status != shim.OK{
		t.Error("failed to delete car")
	}

	// check if car has deleted from worldstate
	carByte,_ = stub.GetState(mockCarID)

	if len(carByte)!=0{
		t.Error("car wasn't deleted from worldstate")
	}
	t.Log("successfully deleted mocked car")
}
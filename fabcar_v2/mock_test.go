package main

import (
	"encoding/json"
	"fmt"
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

	callers = []mockCaller{
		{
			commonName:  "dealer1",
			pemFilePath: "mockData/dealer1_cert",
		},
		{
			pemFilePath: "./mockData/customer1_cert",
			commonName:  "customer1",
		},
		{
			pemFilePath: "./mockData/customer2_cert",
			commonName:  "customer2",
		},
	}
)

func TestCreate(t *testing.T) {
	// create a new mock stub
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// set to caller to dealer
	setCaller(stub, callers[0].readPEM())

	result := stub.MockInvoke("txID",
		[][]byte{
			[]byte("create"),  // function name
			[]byte(mockCarID), // args 0th index
			[]byte(mockCar.Make),
			[]byte(mockCar.Model),
			[]byte(mockCar.Colour),
		})

	if result.Status != shim.OK {
		t.Errorf("failed to create a new car : %s", result.Message)
	}

	// now check worldstate for carID
	carByte, err := stub.GetState(mockCarID)
	if err != nil {
		t.Error(err)
	}
	var car Car
	json.Unmarshal(carByte, &car)
	t.Logf("car from worldstate : %+v", car)

	// log.Println(stub.State) // worldstate dump
}

func TestCreateNonDealer(t *testing.T) {
	// create a new mock stub
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// set caller to a customer
	setCaller(stub, callers[1].readPEM())

	result := stub.MockInvoke("txID",
		[][]byte{
			[]byte("create"),  // function name
			[]byte(mockCarID), // args 0th index
			[]byte(mockCar.Make),
			[]byte(mockCar.Model),
			[]byte(mockCar.Colour),
		})
	if result.Status == shim.OK {
		t.Errorf("non dealer should not be allowed to create new car")
	}
	t.Logf("msg from chaincode : [%s]", result.Message)
}

func TestTransfer(t *testing.T) {
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// insert mock car for transfering the car
	stub.MockTransactionStart("__tx")
	mockCar.Owner = callers[1].commonName
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	newOwner := "newOwner"

	// set caller to customer1
	setCaller(stub, callers[1].readPEM())
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

func TestTransferNotOwned(t *testing.T) {
	stub := shimtest.NewMockStub("__car__", chaincode{})

	// insert mock car for transfering the car
	stub.MockTransactionStart("__tx")
	mockCar.Owner = callers[1].commonName
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	newOwner := "newOwner"

	// set caller to customer2
	setCaller(stub, callers[2].readPEM())
	result := stub.MockInvoke("txID", [][]byte{
		[]byte("transfer"),
		[]byte(mockCarID),
		[]byte(newOwner),
	})

	if result.Status != shim.ERROR {
		t.Error("only owner should be allowed to transfer their cars")
	}

	t.Logf("msg from chiancode : [%s]", result.Message)
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
	mockCar.Owner = callers[1].commonName
	carByte, _ := json.Marshal(mockCar)
	stub.PutState(mockCarID, carByte)
	stub.MockTransactionEnd("__tx")

	// set caller to customer1
	setCaller(stub, callers[1].readPEM())
	result := stub.MockInvoke("txID", [][]byte{
		[]byte("purge"),
		[]byte(mockCarID),
	})

	if result.Status != shim.OK {
		t.Error("failed to delete car")
	}

	// check if car has deleted from worldstate
	carByte, _ = stub.GetState(mockCarID)

	if len(carByte) != 0 {
		t.Error("car wasn't deleted from worldstate")
	}
	t.Log("successfully deleted mocked car")
	t.Log(stub.State) //  worldstate dumb
}

func TestQueryByOwner(t *testing.T) {
	stub := shimtest.NewMockStub("__new__", chaincode{})
	setCaller(stub, callers[0].readPEM()) // set caller to dealer
	// put some mocked cars
	carCount := 5
	for i := 0; i < carCount; i++ {
		stub.MockInvoke("txID",
			[][]byte{
				[]byte("create"),                // function name
				[]byte(fmt.Sprintf("car%d", i)), // args 0th index
				[]byte(fmt.Sprintf("make_%d", i)),
				[]byte(fmt.Sprintf("model_%d", i)),
				[]byte(fmt.Sprintf("colour_%d", i)),
			})
	}

	result := stub.MockInvoke("txID", [][]byte{
		[]byte("queryByOwner"),
		[]byte(callers[0].commonName),
	})

	if result.Status != shim.OK {
		t.Errorf("failed to queryByOwner : %s", result.Message)
	}
	t.Log(string(result.Payload)) // returned query result
}

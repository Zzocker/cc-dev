package contract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// core contains all the core business logic for the chiancode

/*
	- create() : creates a car with owner as client-id
	- transfer(carID,to) : transfer a car owned by client to `to` client
	- query(carId) : returns details about the car
	- queryByOwner(ownerId) : returns array of cars owned by a owner with id `ownerId`
	- purge(carID) : demolish a car after the end of service time

*/

type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

func create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("invalid number of argument, require 1 but got %d", len(args)))
	}
	// args[0] -- marshaled json of car
	var car Car
	err := json.Unmarshal([]byte(args[0]), &car)
	if err != nil {
		return shim.Error(fmt.Sprintf("invalid json object : %s", err.Error()))
	}
	id := stub.GetTxID()
	err = stub.PutState(id, []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(id))
}

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
	ownerID, err := getOwnerID(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	car.Owner = *ownerID
	carByte, _ := json.Marshal(car)
	err = stub.PutState(id, carByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(id))
}

func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 0th index is carID
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("invalid number of argument, require 1 but got %d", len(args)))
	}
	carByte, err := stub.GetState(args[0])
	if err != nil || len(carByte) == 0 {
		return shim.Error(fmt.Sprintf("car with id = %s , not found", args[0]))
	}
	return shim.Success(carByte)
}

func transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : carID , to
	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("invalid number of argument, require 2 but got %d", len(args)))
	}
	carByte, err := stub.GetState(args[0])
	if err != nil || len(carByte) == 0 {
		return shim.Error("car not found")
	}
	var car Car
	json.Unmarshal(carByte, &car)
	ownerID, err := getOwnerID(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if car.Owner != *ownerID {
		return shim.Error("only owner of the car transfer to other owner")
	}
	car.Owner = args[1]
	carByte, _ = json.Marshal(car)
	err = stub.PutState(args[0], carByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func purge(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : carID
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("invalid number of argument, require 1 but got %d", len(args)))
	}
	carByte, err := stub.GetState(args[0])
	if err != nil || len(carByte) == 0 {
		return shim.Error("car not found")
	}
	var car Car
	json.Unmarshal(carByte, &car)
	ownerID, err := getOwnerID(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if car.Owner != *ownerID {
		return shim.Error("only owner of the car can purge")
	}
	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

package contract

import (
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
	return shim.Success(nil)
}

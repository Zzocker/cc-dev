package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type chaincode struct{}

var (
	supportedMethod = map[string]func(stub shim.ChaincodeStubInterface, args []string) peer.Response{
		"create":       create,
		"transfer":     transfer,
		"query":        query,
		"queryByOwner": queryByOwner,
		"purge":        purge,
	}
)

func (chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	methodName, args := stub.GetFunctionAndParameters()
	method, ok := supportedMethod[methodName]
	if !ok {
		return shim.Error(fmt.Sprintf(`method "<%s>" not supported`, methodName))
	}
	return method(stub, args)
}

/*
	fabcar_v2 : use golevel db as worldstate
		1. create(carID,make,model,colour) // only client with IS_DEALER set to true can create a new car
		2. transfer(carID,newOwner) // only owner can transfer their car
		3. query(carID)
		4. queryByOwner(owner) : use composite key
		5. purge(carID) // only owner can purge their car
*/
const IS_DEALER_ATTR = "IS_DEALER"
const OWNER_CAR_INDEX = "ownerID~carID"

var NULL_VALUE = []byte{0x00} // nul value to be used as value for Composite key

type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

func create(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// caller is dealer or not
	identity, err := cid.New(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = identity.AssertAttributeValue(IS_DEALER_ATTR, "true")
	if err != nil {
		return shim.Error("only dealer can create a new car : " + err.Error())
	}

	// args : carID , make , model , colour
	if len(args) != 4 {
		return invalidNumberOfArgument(args, 4)
	}
	// carID should be already existing
	carByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// size of carByte will be zero for non-existing state
	if len(carByte) != 0 {
		return shim.Error(fmt.Sprintf("car %s already exists", args[0]))
	}
	ownerID := getCallerID(stub)
	car := Car{
		Make:   args[1],
		Model:  args[2],
		Colour: args[3],
		Owner:  ownerID, // set ownership of car to calling client
	}
	carByte, _ = json.Marshal(car)

	// put car state into worldstate
	// with id = carID
	// value = carByte
	err = stub.PutState(args[0], carByte)
	if err != nil {
		return shim.Error(err.Error())
	}

	// new composite key with attributes : ownerID, carID
	key, err := stub.CreateCompositeKey(OWNER_CAR_INDEX, []string{ownerID, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, NULL_VALUE)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func transfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : carID , newOwner
	if len(args) != 2 {
		return invalidNumberOfArgument(args, 2)
	}
	// get older version of state
	carByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(carByte) == 0 {
		return shim.Error("car not found")
	}
	var car Car
	json.Unmarshal(carByte, &car)

	// ownership check
	ownerID := getCallerID(stub)
	if car.Owner != ownerID {
		return shim.Error("only owner can transfer car to another")
	}

	// set owner to new owner (args[1])
	car.Owner = args[1]

	carByte, _ = json.Marshal(car)
	// update the state of car
	err = stub.PutState(args[0], carByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : carID
	if len(args) != 1 {
		return invalidNumberOfArgument(args, 1)
	}
	carByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(carByte) == 0 {
		return shim.Error("car not found")
	}
	return shim.Success(carByte)
}

func queryByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : ownerID
	if len(args) != 1 {
		return invalidNumberOfArgument(args, 1)
	}

	// returns all composite key starting pattern ownerID*
	iterator, err := stub.GetStateByPartialCompositeKey(OWNER_CAR_INDEX, []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	// close iterator to free up the resources
	defer iterator.Close()
	var cars []Car

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			continue
		}
		// kv.Key = ownerID~carIDdealer1car0 , need to split and get carID
		// SplitCompositeKey : returns index:string , attributes:[]string , error
		_, attr, _ := stub.SplitCompositeKey(kv.Key)
		carID := attr[1]

		// get car with returnred carID
		carByte, _ := stub.GetState(carID)
		var car Car
		json.Unmarshal(carByte, &car)
		cars = append(cars, car)
	}
	out, _ := json.Marshal(cars)
	return shim.Success(out)
}

func purge(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// args : carID
	if len(args) != 1 {
		return invalidNumberOfArgument(args, 1)
	}
	carByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(carByte) == 0 {
		// car doesn't exists
		return shim.Error("car not found")
	}
	var car Car
	json.Unmarshal(carByte, &car)

	// ownership check
	ownerID := getCallerID(stub)
	if car.Owner != ownerID {
		return shim.Error("only owner can purge their car")
	}

	// delete car state
	err = stub.DelState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	// delete compiste key also
	key, err := stub.CreateCompositeKey(OWNER_CAR_INDEX, []string{ownerID, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func main() {
	if err := shim.Start(chaincode{}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// util functions
func invalidNumberOfArgument(args []string, wanted int) peer.Response {
	return shim.Error(fmt.Sprintf("invalid number of argument , wanted : %d , but got %d", len(args), wanted))
}

package contract

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

func TestCreate(t *testing.T) {
	// create mock stub with given chaincode
	stub := shimtest.NewMockStub("_car_", new(Chaincode))

	myCar := Car{
		Make:   "Toyota",
		Model:  "Prius",
		Colour: "blue",
		Owner:  "Tomoko",
	}
	carByte, _ := json.Marshal(myCar)

	// make mock invoke with function "create"
	// input : <Tx ID > , [][]byte{} - array of argument with 0th index as method name
	txId := "tx1"
	result := stub.MockInvoke(txId, [][]byte{[]byte("create"), carByte})

	// result is a protocol buffer message with fields
	// - Message : string ,
	// - Status : int,
	// Payload : []byte
	if result.Status != shim.OK {
		t.Errorf("failed to invoke create a new car : %s", result.Message)
	}

	// check if worldstate is updated or not
	carByte, err := stub.GetState(txId)
	if err != nil || len(carByte) == 0 {
		t.Error("car is not created")
	}

	var wantCar Car
	err = json.Unmarshal(carByte, &wantCar)
	if err != nil {
		t.Errorf("data formate miss-match %s", err.Error())
	}

	// check if input car and car from worldstate are same
	if !reflect.DeepEqual(myCar, wantCar) {
		t.Errorf("wanted : %+v , but got : %+v", myCar, wantCar)
	}
}

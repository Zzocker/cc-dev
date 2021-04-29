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

func TestQuery(t *testing.T) {
	stub := shimtest.NewMockStub("_car_", new(Chaincode))

	t.Run("found", func(t *testing.T) {
		// put mock data into worldstate
		myCar := Car{
			Make:   "Toyota",
			Model:  "Prius",
			Colour: "blue",
			Owner:  "Tomoko",
		}
		carByte, _ := json.Marshal(myCar)
		id := "car1"
		// update worldstate require tx to be started and ended
		stub.MockTransactionStart("__query")
		stub.PutState("car1", carByte)
		stub.MockTransactionEnd("__query")
		//
		gotCarByte, err := stub.GetState(id)
		if err != nil || len(gotCarByte) == 0 {
			t.Errorf("query should return car with id = %s , but got empty data", id)
		}
		var gotCar Car
		err = json.Unmarshal(gotCarByte, &gotCar)
		if err != nil {
			t.Error("car data formate miss-match")
		}
		if !reflect.DeepEqual(myCar, gotCar) {
			t.Errorf("wantted : %+v, but got : %+v", myCar, gotCar)
		}
	})

	t.Run("not-found", func(t *testing.T) {
		result := stub.MockInvoke("__query__", [][]byte{[]byte("not_found")})
		if result.Status != shim.ERROR {
			t.Errorf("should return not found error , but got otherwise %s", string(result.Payload))
		}
	})
}
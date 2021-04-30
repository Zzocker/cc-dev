package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func getOwnerID(stub shim.ChaincodeStubInterface) (*string, error) {
	identity, err := cid.New(stub)
	if err != nil {
		return nil, err
	}
	mspID, err := identity.GetMSPID()
	if err != nil {
		return nil, err
	}
	cert, err := identity.GetX509Certificate()
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("%s_%s", mspID, cert.Subject.CommonName)
	return &id, nil
}

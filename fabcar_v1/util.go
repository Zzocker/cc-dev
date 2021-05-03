package main

import (
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// returns read certificate of caller client and return value commonName field
func getCallerID(stub shim.ChaincodeStubInterface) string {
	identity, _ := cid.New(stub)
	cert, _ := identity.GetX509Certificate()
	return cert.Subject.CommonName
}

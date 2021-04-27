package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/Zzocker/cc-dev/contract"
)

func main() {
	cc := new(contract.Chaincode)
	if err := shim.Start(cc); err != nil {
		fmt.Printf("failed to start the chaincode : %v", err)
		os.Exit(1)
	}
}

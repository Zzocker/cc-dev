package contract

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-protos-go/msp"
)

type mockCreator struct {
	name    string
	certPem string
}

var (
	testMSPID = "DevMSP"
	user1     = mockCreator{
		name: "user1",
		certPem: `-----BEGIN CERTIFICATE-----
MIICgzCCAimgAwIBAgIUUqD8w0a2/3P2N+zm1uvZfw8PUSMwCgYIKoZIzj0EAwIw
ZzELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNh
biBGcmFuY2lzY28xEzARBgNVBAoTCmRldm9yZy5jb20xFjAUBgNVBAMTDWNhLmRl
dm9yZy5jb20wHhcNMjEwNDI0MTEyMjAwWhcNMjIwNDI0MTEyNzAwWjBCMTAwDQYD
VQQLEwZjbGllbnQwCwYDVQQLEwRvcmcxMBIGA1UECxMLZGVwYXJ0bWVudDExDjAM
BgNVBAMTBXVzZXIxMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPtYJrextCxlB
Ot6ENC9BRFPPVD8OotBq0peo0cOg32GEnyUrg45KBIMNYibResW1M3WmoSFt67LI
nyQLFTfA6qOB1zCB1DAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADAdBgNV
HQ4EFgQUMYaqrhJJ03f/LXWRkMSgGzlcOxYwKwYDVR0jBCQwIoAg7IbBCX1GeUPY
ATwsM/QGL+vfNXlu0sWLgrYk4IsDvHQwaAYIKgMEBQYHCAEEXHsiYXR0cnMiOnsi
aGYuQWZmaWxpYXRpb24iOiJvcmcxLmRlcGFydG1lbnQxIiwiaGYuRW5yb2xsbWVu
dElEIjoidXNlcjEiLCJoZi5UeXBlIjoiY2xpZW50In19MAoGCCqGSM49BAMCA0gA
MEUCIQDVhLKdIvu27RqjNHP099RzFNSHg404o9e2eF9pQSM+nAIgFQGNrKCHDGAU
jlHe30Y2CrQAH5b5FuX8uSadD6NGzgQ=
-----END CERTIFICATE-----`,
	}
)

func TestGetOwnerID(t *testing.T) {
	stub := shimtest.NewMockStub("_-", new(Chaincode))

	setCreator(stub, user1.certPem)

	wantID := fmt.Sprintf("%s_%s", testMSPID, user1.name)

	gotID, err := getOwnerID(stub)
	if err != nil {
		t.Error(err)
	}
	if *gotID != wantID {
		t.Errorf("ownerID missmatch wanted = %s , but got = %s", wantID, *gotID)
	}
}

func setCreator(stub *shimtest.MockStub, pem string) {
	identity := msp.SerializedIdentity{
		Mspid:   "DevMSP",
		IdBytes: []byte(pem),
	}
	d, _ := proto.Marshal(&identity)
	stub.Creator = d
}

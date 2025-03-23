package config

import (
	"crypto/ecdsa"
	"robotech/utils/ethUtils"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func Test_initEnv(t *testing.T) {
	t.Logf("\nInitEnv testing...")

	client := ethUtils.GetClient()
	if client == nil {
		t.Fatalf("client is nil")
	}

	publicKey := sKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatalf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Logf("\nfromAddress: %v", fromAddress)
}

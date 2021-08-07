package lendingpoolevent_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func Test_Topic0(t *testing.T) {
	AssetIndexUpdated := crypto.Keccak256Hash([]byte("AssetIndexUpdated(address,uint256)"))
	fmt.Println(AssetIndexUpdated)
	if AssetIndexUpdated.Hex() != "0x5777ca300dfe5bead41006fbce4389794dbc0ed8d6cccebfaf94630aa04184bc" {
		t.Fatalf("expect: 0x5777ca300dfe5bead41006fbce4389794dbc0ed8d6cccebfaf94630aa04184bc\ngot   : %v", AssetIndexUpdated)
	}
}

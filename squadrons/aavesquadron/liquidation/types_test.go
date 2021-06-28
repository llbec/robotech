package liquidation

import (
	"fmt"
	"testing"
)

func Test_Generateconfig(t *testing.T) {
	err := GeneratConfig()
	if err != nil {
		t.Fatal(err)
	}
	initEnv()
	fmt.Println(Reserves)
	fmt.Println(Contracts)
	fmt.Println(Clients)
}

func Test_GetAuth(t *testing.T) {
	initEnv()
	fmt.Println(GetAuther())
}

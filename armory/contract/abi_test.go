package contract

import (
	"fmt"
	"testing"
)

func Test_getFunction(t *testing.T) {
	//funStr := "function getReserveData(address asset) external view returns (DataTypes.ReserveData memory);"
	funStr := "function getReserveData(address asset) external view returns (uint256);"
	str, err := getFunction(funStr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(str)
}

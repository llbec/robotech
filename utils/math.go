package utils

import (
	"math"
	"math/big"
)

func BigToRecognizable(v *big.Int, decimal int) *big.Float {
	z := big.NewFloat(1)
	return z.Quo(big.NewFloat(1).SetInt(v), big.NewFloat(math.Pow10(int(decimal))))
}

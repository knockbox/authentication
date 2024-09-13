package utils

import (
	"crypto/rand"
	"math/big"
)

// CryptoRandom returns an int in the range [0, max)
func CryptoRandom(max int) (int64, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return bigInt.Int64(), nil
}

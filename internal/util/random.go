package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GenerateRandomDigits(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("n must be a positive integer")
	}

	builder := strings.Builder{}
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("error generating random number: %w", err)
		}
		builder.WriteString(num.String())
	}

	return builder.String(), nil
}

func RandomNumber(start, end int) int {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(end-start)))
	if err != nil {
		return 0
	}
	return int(num.Int64()) + start
}

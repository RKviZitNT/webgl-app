package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
)

func GenerateRandomCode(lenght int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, lenght)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

func ReadStruct(data interface{}, target interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

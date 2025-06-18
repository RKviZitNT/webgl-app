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

func ParseInterfaceToJSON(msgData interface{}, output interface{}) error {
	bytes, err := json.Marshal(msgData)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, output)
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ParseStringToJSON(data string, output interface{}) error {
	bytes := []byte(data)
	return json.Unmarshal(bytes, output)
}

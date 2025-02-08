package helper

import (
	"crypto/rand"
)

func RandomNumbers(length int) (string, error) {
	const numbers = "1234567890"

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	numLenght := len(numbers)

	for i := 0; i < length; i++ {
		buffer[i] = numbers[int(buffer[i])%numLenght]
	}

	return string(buffer), nil
}

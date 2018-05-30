package main

import (
	"crypto/rand"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomBytes2(n int) []byte {
	return make([]byte, n)
}

func reverseInplace(buffer *[]byte) {
	tmp := *buffer

	j := len(tmp) - 1
	i := 0
	for i < j {
		t := tmp[j]
		tmp[j] = tmp[i]
		tmp[i] = t

		i++
		j--
	}

	buffer = &tmp
}

func reverseInplace2(buffer *[]byte) {
	tmp := *buffer
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	buffer = &tmp
}

func abs(n int16) int16 {
	if n < 0 {
		return -n
	}
	return n
}

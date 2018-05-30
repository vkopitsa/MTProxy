package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseInplace(t *testing.T) {
	data := []byte{171, 47, 123, 204}
	expect := []byte{204, 123, 47, 171}
	reverseInplace(&data)
	assert.Equal(t, expect, data)
}

func BenchmarkReverseInplace(b *testing.B) {
	data := []byte{171, 47, 123, 204}

	for n := 0; n < b.N; n++ {
		reverseInplace(&data)
	}
}

func TestReverseInplace2(t *testing.T) {
	data := []byte{171, 47, 123, 204}
	expect := []byte{204, 123, 47, 171}
	reverseInplace2(&data)
	assert.Equal(t, expect, data)
}

func BenchmarkReverseInplace2(b *testing.B) {
	data := []byte{171, 47, 123, 204}

	for n := 0; n < b.N; n++ {
		reverseInplace2(&data)
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	buf64, err := GenerateRandomBytes(64)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(buf64))
}

func BenchmarkGenerateRandomBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateRandomBytes(64)
	}
}

func TestGenerateRandomBytes2(t *testing.T) {
	buf64 := GenerateRandomBytes2(64)
	assert.Equal(t, 64, len(buf64))
}

func BenchmarkGenerateRandomBytes2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateRandomBytes2(64)
	}
}

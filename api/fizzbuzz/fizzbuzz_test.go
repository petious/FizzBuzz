package fizzbuzz

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGenerateFizzBuzz(t *testing.T) {
	res := GenerateFizzBuzz(3, 5, 15, "fizz", "buzz")
	assert.Equal(t, "12fizz4buzzfizz78fizzbuzz11fizz1314fizzbuzz", res)
}

func TestGetCorrectStringForNumber(t *testing.T) {
	res := GetCorrectStringForNumber(1, 1, 2, "fizz", "buzz")
	assert.Equal(t, "fizz", res)

	res = GetCorrectStringForNumber(2, 1, 2, "fizz", "buzz")
	assert.Equal(t, "fizzbuzz", res)

	res = GetCorrectStringForNumber(1, 3, 2, "fizz", "buzz")
	assert.Equal(t, "1", res)
}

package fizzbuzz

import (
	"math"
	"strconv"
)

func GenerateFizzBuzz(int1 int, int2 int, limit int, str1 string, str2 string) string {
	resultString := ""
	for i := 1; i <= limit; i++ {
		resultString += GetCorrectStringForNumber(i, int1, int2, str1, str2)
	}
	return resultString
}

func GetCorrectStringForNumber(currentNumber int, int1 int, int2 int, str1 string, str2 string) string {

	remainerInt1 := math.Mod(float64(currentNumber), float64(int1))
	remainerInt2 := math.Mod(float64(currentNumber), float64(int2))

	if remainerInt1 == 0 && remainerInt2 == 0 {
		return str1 + str2
	}
	if remainerInt1 == 0 {
		return str1
	}

	if remainerInt2 == 0 {
		return str2
	}

	return strconv.Itoa(currentNumber)
}

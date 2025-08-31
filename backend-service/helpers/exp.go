package helpers

import (
	"fmt"
	"math"
)

func CalculateExpNeeded(level int) int {
	cnf := NewConfig()
	base := cnf.GetFloat64("BASE_EXP")
	factor := cnf.GetFloat64("EXP_FACTOR")
	return int(base * math.Pow(float64(level), factor))
}

func CheckExpGain(difficulty string) (int, error) {
	switch difficulty {
	case "easy":
		return 50, nil
	case "normal":
		return 100, nil
	case "hard":
		return 200, nil
	default:
		return 0, fmt.Errorf("Difficulty invalid")
	}
}

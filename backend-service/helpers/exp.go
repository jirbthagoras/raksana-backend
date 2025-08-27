package helpers

import "math"

func CalculateExpNeeded(level int) int {
	cnf := NewConfig()
	base := cnf.GetFloat64("BASE_EXP")
	factor := cnf.GetFloat64("EXP_FACTOR")
	return int(base * math.Pow(float64(level), factor))
}

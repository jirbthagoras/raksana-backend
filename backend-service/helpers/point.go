package helpers

var LevelMultipliers = map[int]float64{
	1: 1.0, 2: 1.0, 3: 1.0,
	4: 1.2, 5: 1.2, 6: 1.2,
	7: 1.5, 8: 1.5, 9: 1.5,
	10: 1.8, 11: 1.8, 12: 1.8,
	13: 2.0, 14: 2.0, 15: 2.0,
}

func GetMultiplier(level int, streak int) float64 {
	multiplier := LevelMultipliers[level]
	streakBonus := float64(streak/5) * 0.1
	return multiplier + streakBonus
}

package helpers

import (
	"jirbthagoras/raksana-backend/repositories"
	"math/rand"
)

func WeightedRandomPick(habits []repositories.Habit) repositories.Habit {
	totalWeight := 0
	for _, h := range habits {
		totalWeight += int(h.Weight)
	}

	r := rand.Intn(totalWeight) + 1

	for _, h := range habits {
		r -= int(h.Weight)
		if r <= 0 {
			return h
		}
	}

	return habits[0]
}

func PickMultiple(habits []repositories.Habit, count int) []repositories.Habit {
	if count > len(habits) {
		count = len(habits)
	}

	chosen := []repositories.Habit{}
	used := map[int]bool{}

	for len(chosen) < count {
		h := WeightedRandomPick(habits)
		if !used[int(h.ID)] {
			chosen = append(chosen, h)
			used[int(h.ID)] = true
		}
	}

	return chosen
}

package utils

import "math/rand"

func IsError() bool {
	// throw errors 1 time out of 100
	if rand.Intn(100) < 99 {
		return false
	}
	return true
}

func IsErrorMoreLikely() bool {
	// throw errors 10 times out of 100
	if rand.Intn(100) < 90 {
		return false
	}
	return true
}

func IsErrorPrettyLikely() bool {
	// throw errors 50 times out of 100
	if rand.Intn(100) < 50 {
		return false
	}
	return true
}


func IsErrorVeryLikely() bool {
	// throw errors 80 times out of 100
	if rand.Intn(100) < 20 {
		return false
	}
	return true
}

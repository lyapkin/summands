package main

func upperBound(target int, length int) int {
	var sum int
	for i := 1; i < length; i++ {
		sum += i
	}
	return target - sum
}
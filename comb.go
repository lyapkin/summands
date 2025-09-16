package main

import (
	"fmt"
)


type Combs interface {
	FindCombs() (*[][]int, int)
}

func NewCombs(target int, length int, ub int, dir string, updateState func(int)) Combs {
	if ub == 0 {
		ub = upperBound(target, length)
	}

	file := NewFile(fmt.Sprintf("%d-%d-%d", target, length, ub), dir)

	return &_Combs{
		target: target,
		length: length,
		ub: ub,
		file: file,
		updateState: updateState,
	}
}

type _Combs struct {
	target int
	length int
	ub int
	file File
	updateState func (int)
	found int
}

func (c *_Combs) FindCombs() (*[][]int, int) {
	results := make([][]int, 0, BATCH)
	defer func() {
		c.file.Close()
		c.file = nil
	}()

	c.calculate(make([]int, 0, c.length), c.target, c.length, 1, &results, c.ub)
	c.file.Write(&results)
	c.file.Save()
	
	return &results, c.found + len(results)
}

func (c *_Combs) calculate(curCombination []int, target int, length int, start int, results *[][]int, ub int) {
	if target == 0 && length == 0 {
		newComb := make([]int, len(curCombination), cap(curCombination))
		copy(newComb, curCombination)
		*results = append(*results, newComb)
		if len(*results) == cap(*results) {
			c.file.Write(results)
			c.found += len(*results)
			*results = (*results)[:0]

			c.updateState(c.found)
		}
		return
	}

	if length <= 0 {
		return
	}

	for i := start; i <= min(target, ub); i++ {
		if i * length > target {
			break
		}

		if length == 1 {
			val := min(target, ub)
			c.calculate(append(curCombination, val), target - val, 0, val + 1, results, ub)
			break
		}

		c.calculate(append(curCombination, i), target - i, length - 1, i + 1, results, ub)
	}
}
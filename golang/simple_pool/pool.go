package main

import (
	"github.com/vearne/golib/utils"
	"log"
	"strconv"
)

func Judge(key interface{}) *utils.GPResult {
	result := &utils.GPResult{}
	num, _ := strconv.Atoi(key.(string))
	if num < 450 {
		result.Value = true
	} else {
		result.Value = false
	}
	return result
}

func main() {
	p := utils.NewGPool(30)

	slice := make([]interface{}, 0)
	for i := 0; i < 1000; i++ {
		slice = append(slice, strconv.Itoa(i))
	}

	result := make([]bool, 0, 10)
	trueCount := 0
	falseCount := 0
	for item := range p.ApplyAsync(Judge, slice) {
		value := item.Value.(bool)
		result = append(result, value)
		if value {
			trueCount++
		} else {
			falseCount++
		}
	}

	log.Printf("cancel, %v, true:%v, false:%v\n", len(result), trueCount, falseCount)
}

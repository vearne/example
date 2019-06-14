package main

import (
	"context"
	"fmt"
	"github.com/vearne/golib/utils"
	"log"
	"strconv"
	"time"
)



func JudgeStrWithContext2(ctx context.Context, key interface{}) *utils.GPResult {
	num, _ := strconv.Atoi(key.(string))
	result := &utils.GPResult{}

	var canceled bool = false

	for i := 0; i < 60; i++ {
		select {
		case <-ctx.Done():
			canceled = true
			result.Value = false
			result.Err = fmt.Errorf("normal termination")
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

	if !canceled {
		if num < 450 {
			result.Value = true
		} else {
			result.Value = false
		}
	}

	return result
}

func main() {
	cxt, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	p := utils.NewGContextPool(cxt,30)

	slice := make([]interface{}, 0)
	for i := 0; i < 1000; i++ {
		slice = append(slice, strconv.Itoa(i))
	}


	result := make([]*utils.GPResult, 0, 10)
	trueCount := 0
	falseCount := 0

	start := time.Now()
	for item := range p.ApplyAsync(JudgeStrWithContext2, slice) {
		result = append(result, item)
		if item.Err!= nil{
			//log.Println("cancel", item.Err)
			continue
		}
		if item.Value.(bool) {
			trueCount++
		} else {
			falseCount++
		}
	}

	log.Printf("cancel, %v, true:%v, false:%v, cost:%v\n", len(result),
		trueCount, falseCount, time.Since(start))
}

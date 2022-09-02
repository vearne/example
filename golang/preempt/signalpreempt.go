package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	var n int
	go func() {
		for {
			n++
			fmt.Println("goroutine1:", n)
		}
	}()
	//runtime.Gosched()
	for {
		n++
		fmt.Println("goroutine2:", n)
	}
}

func fibR(i int) int {
	if i < 2 {
		return i
	}
	return fibR(i-1) + fibR(i-2)
}

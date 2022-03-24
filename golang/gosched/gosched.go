package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 100; i++ {
		fmt.Println("Goroutine1:", i)
		//time.Sleep(500 * time.Millisecond)
		if i == 5 {
			go func() {
				for i := 0; i < 100; i++ {
					fmt.Println("Goroutine2:", i)
					//time.Sleep(500 * time.Millisecond)
				}
			}()
			runtime.Gosched()
		}
	}
}

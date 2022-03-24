package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(1)
	var n int
	go func() {
		for {
			n++
			//fmt.Println(n)
		}
	}()
	//runtime.Gosched()
	for {
		n++
	}
}

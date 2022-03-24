package main

import (
	"math"
	"os"
	"runtime"
	"runtime/trace"
)

func main() {
	runtime.GOMAXPROCS(1)
	//创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	for i := 0; i < math.MaxInt32; i++ {
		//fmt.Println("Goroutine1:", i)
		if i == 5 {
			go func() {
				for i := 0; i < math.MaxInt32; i++ {
					//fmt.Println("Goroutine2:", i)
				}
			}()
			runtime.Gosched()
		}
	}
}

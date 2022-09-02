package main

import (
	_ "embed"
)

//go:embed variable/data.txt
var data1 string

//go:embed variable/data2.txt
var data2 []byte

func main() {
	println("data1:", data1)
	println("data2:", string(data2))
}

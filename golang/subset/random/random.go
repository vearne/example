package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

/*
	data <- read.table("/tmp/datafile.csv",header=TRUE, sep=",")
	ggplot(data, aes(X, Y)) + geom_bar(stat = 'identity')
 */
// 假定 50 client, 100 backend
const (
	clientSize  = 50
	backendSize = 100
	subsetSize  = 50
	// 假定Client的负载是均衡的，某个时间段里每个Client处理100000个请求
	clientTotalReq = 100000
)

func init() {
	rand.Seed(0)
}

func main() {
	/*
		Client的编号从 0 到 clientSize -1
		Backend的编号从 0 到 backendSize - 1
	*/
	backends := make([]int, 0)
	for i := 0; i < backendSize; i++ {
		backends = append(backends, i)
	}

	// clientID -> client上对应的负载
	var loadMap map[int]int = make(map[int]int)

	for i := 0; i < clientSize; i++ {
		subset := GetRandomSubSet(backends, subsetSize)
		//fmt.Println("len(subset)", len(subset))
		for _, backendID := range subset {
			loadMap[backendID] += clientTotalReq / subsetSize
		}
	}

	expected := clientSize * clientTotalReq / backendSize
	maxMargin := 0

	// print result
	for i := 1; i < backendSize+1; i++ {
		margin := int(math.Abs(float64(loadMap[i] - expected)))
		fmt.Printf("backends[%d], load: %d, margin: %d\n", i, loadMap[i], margin)
		if margin > maxMargin {
			maxMargin = margin
		}
	}

	fmt.Printf("expected:%d, maxMargin:%d\n", expected, maxMargin)
	// 写文件
	WriteCSV(loadMap)
}

func GetRandomSubSet(backends []int, subsetSize int) []int {
	ans := make([]int, 0)
	for i := 0; i < subsetSize; i++ {
		ans = append(ans, backends[rand.Intn(len(backends))])
	}
	return ans
}


func WriteCSV(loadMap map[int]int){
	// 创建一个 tutorials.csv 文件
	csvFile, err := os.Create("/tmp/datafile.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Write([]string{"X", "Y"})
	for i:=0;i<len(loadMap);i++{
		writer.Write([]string{"clientID:" + strconv.Itoa(i), strconv.Itoa(loadMap[i])})
	}
	writer.Flush()
}
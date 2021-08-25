package main

import (
	"encoding/csv"
	"fmt"
	"github.com/vearne/golib/utils"
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
	//clientSize  = 50
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


	// clientID -> client上连接池的变化量
	var connMap map[int]int = make(map[int]int)

	for i := 0; i < clientSize; i++ {
		backends := make([]int, 0)
		for i := 0; i < backendSize; i++ {
			backends = append(backends, i)
		}
		subsetOld := Subset(backends, i, subsetSize)
		backends = make([]int, 0)
		for i := 0; i < backendSize; i++ {
			// 假定 backend-10 宕机
			if i != 10{
				backends = append(backends, i)
			}

		}
		subsetNew := Subset(backends, i, subsetSize)
		connOld := utils.NewIntSet()
		connOld.AddAll(subsetOld)

		connNew := utils.NewIntSet()
		connNew.AddAll(subsetNew)
		connNew.RemoveAll(connOld)
		connMap[i] = connNew.Size()
		fmt.Printf("Client[%d], 连接变化量: %d\n", i, connMap[i])
	}
	totalConn := clientSize * subsetSize
	totalChange := 0
	for _, changeCount := range connMap{
		totalChange += changeCount
	}

	fmt.Printf("所有连接数:%d, 重建的连接总数:%d, 百分比: %f %%\n", totalConn, totalChange, float64(totalChange)/float64(totalConn) * 100)
	// 写文件
	WriteCSV(connMap)
}

// 采用的方案是<<SRE Google运维解密>>书中提到的子集选择算法二：确定性子集。
// clientId是将原client的Ip地址做了CRC处理，转成int
// subsetSize就是子集的大小
func Subset(backends []int, clientID, subsetSize int) []int {

	subsetCount := len(backends) / subsetSize

	// Group clients into rounds; each round uses the same shuffled list:
	round := clientID / subsetCount

	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(backends), func(i, j int) { backends[i], backends[j] = backends[j], backends[i] })

	// The subset id corresponding to the current client:
	subsetID := clientID % subsetCount

	start := subsetID * subsetSize

	//fmt.Printf("clientID:%v, round:%v, start:%v, end:%v\n", clientID, round, start, start+subsetSize)
	return backends[start : start+subsetSize]
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
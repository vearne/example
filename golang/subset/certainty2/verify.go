package main

import (
	"encoding/csv"
	"fmt"
	"github.com/spaolacci/murmur3"
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
	//backends := make([]int, 0)
	//for i := 0; i < backendSize; i++ {
	//	backends = append(backends, i)
	//}

	// clientID -> client上对应的负载
	var loadMap map[int]int = make(map[int]int)



	for i := 0; i < clientSize; i++ {
		backends := make([]int, 0)
		for i := 0; i < backendSize; i++ {
			backends = append(backends, i)
		}

		// 构造ClientID
		// 这里根据IP通过hash计算ClientID
		//IP的范围从 192.168.1.1 ~ 192.168.1.50
		clientID := calcuMurmur64(fmt.Sprintf("192.168.1.%d", i+1))
		subset := Subset(backends, clientID, subsetSize)
		//fmt.Println("len(subset)", len(subset))
		for _, backendID := range subset {
			loadMap[backendID] += clientTotalReq / len(subset)
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
	return backends[start : start+subsetSize]
}

func WriteCSV(loadMap map[int]int) {
	// 创建一个 tutorials.csv 文件
	csvFile, err := os.Create("/tmp/datafile.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	writer.Write([]string{"X", "Y"})
	for i := 0; i < len(loadMap); i++ {
		writer.Write([]string{"clientID:" + strconv.Itoa(i), strconv.Itoa(loadMap[i])})
	}
	writer.Flush()
}

func calcuMurmur64(str string) int {
	h32 := murmur3.New64()
	h32.Write([]byte(str))
	x := int(h32.Sum64())
	if x < 0 {
		return -x
	}
	return x
}

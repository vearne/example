package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"github.com/vearne/golib/utils"
	"hash/crc32"
	"math"
	"math/big"
	"math/rand"
	"net"
)

// 假定 50 client, 100 backend
const (
	clientSize = 50
	backendSize = 240
	subsetSize = 50
	clientTotalReq= 100000
)



func test(){

	// 假定IP的范围是 192.168.1.1 ~ 192.168.1.255之间
	var weightMap map[int]int = make(map[int]int)

	for i:=0;i<clientSize;i++{
		//if i == 15||i == 16{
		//	continue
		//}
		ip := fmt.Sprintf("192.168.1.%d", i)
		// init backend
		backends:= make([]int, backendSize)
		for i:=0;i<backendSize;i++{
			backends[i] = i+1
		}
		//clientId := calcuCrc32(ip)
		//clientId := i + 53
		//clientId := calcuMurmur32(ip)
		//clientId := rand.Intn(100)

		clientId := int(InetAtoN(ip))
		fmt.Println("clientId:", clientId)

		setSetSlice := subset(backends, clientId, subsetSize)
		for _, backendId:=range setSetSlice{
			weightMap[backendId]++
		}
	}
	expected := clientSize*subsetSize/backendSize

	maxMargin := 0

	// print result
	for i:=1;i<backendSize+1;i++{
		margin := int(math.Abs(float64(weightMap[i] - expected)))
		fmt.Printf("backends[%d], %d, %d\n", i, weightMap[i], margin)
		if margin > maxMargin{
			maxMargin = margin
		}
	}

	fmt.Printf("expected:%d, maxMargin:%d\n", clientSize*subsetSize/backendSize, maxMargin)
}

func test2(){
	// 假定IP的范围是 192.168.1.1 ~ 192.168.1.255之间
	expected := clientSize*subsetSize/backendSize
	maxMargin := 0

	for i:=0;i<clientSize;i++{
		if i == 15||i == 16{
			continue
		}
		ip := fmt.Sprintf("192.168.1.%d", i)
		// init backend
		backends:= make([]int, 0, backendSize)
		for i:=1;i<backendSize+1;i++{
			backends = append(backends, i)
		}

		backends2:= make([]int, 0, backendSize)
		for i:=1;i<backendSize+1;i++{
			if i == 2 {
				continue
			}
			backends2 = append(backends2, i)
		}
		//clientId := calcuCrc32(ip)
		//clientId := i + 53
		//clientId := calcuMurmur32(ip)
		//clientId := rand.Intn(100)

		clientId := int(InetAtoN(ip))
		fmt.Println("clientId:", clientId)

		setSetSlice := subset(backends, clientId, subsetSize)
		setSetSlice2 := subset(backends2, clientId, subsetSize)
		set1 := utils.NewIntSet()
		set1.AddAll(setSetSlice)

		set2 := utils.NewIntSet()
		set2.AddAll(setSetSlice2)


		set2.RemoveAll(set1)
		margin := set2.Size()
		fmt.Println(set2.Size(), set1.Size())

		if margin > maxMargin{
			maxMargin = margin
		}
	}
	fmt.Println("expected", expected, "maxMargin", maxMargin)
}

func test3(){

	// 假定IP的范围是 192.168.1.1 ~ 192.168.1.255之间
	var weightMap map[int]int = make(map[int]int)

	ipList := []string{
		"10.39.141.25",
		"10.39.141.26",
		"10.39.141.27",
		"10.39.141.28",
		"10.39.141.29",
		"10.39.141.31",
		"10.39.141.33",
		"10.39.141.35",
		"10.39.141.37",
		"10.39.141.39",
		"10.39.141.40",
		"10.39.141.41",
		"10.39.141.191",
		"10.39.140.192",
		"10.39.140.193",
		"10.39.140.194",
		"10.39.140.195",
		"10.39.140.196",
		"10.39.140.197",
		"10.39.140.198",
		"10.39.140.199",
		"10.39.140.200",
		"10.39.140.201",
		"10.39.140.202",
		"10.39.140.203",
		"10.39.140.204",
		"10.39.140.205",
		"10.39.140.206",
		"10.39.140.207",
		"10.39.140.208",
		"10.39.140.209",
		"10.39.140.210",
		"10.39.140.211",
		"10.39.140.212",
		"10.39.140.213",
		"10.39.140.214",
		"10.39.140.215",
		"10.39.140.216",
		"10.39.140.217",
		"10.39.140.218",
	}

	for i:=0;i<len(ipList);i++{
		ip := ipList[i]
		// init backend
		backends:= make([]int, backendSize)
		for i:=0;i<backendSize;i++{
			backends[i] = i+1
		}
		//clientId := calcuCrc32(ip)
		//clientId := i + 53
		//clientId := calcuMurmur32(ip)
		//clientId := rand.Intn(100)

		clientId := int(InetAtoN(ip))
		fmt.Println("clientId:", clientId)

		setSetSlice := subset(backends, clientId, subsetSize)
		for _, backendId:=range setSetSlice{
			weightMap[backendId]++
		}
	}
	expected := clientSize*subsetSize/backendSize

	maxMargin := 0

	// print result
	for i:=1;i<backendSize+1;i++{
		margin := int(math.Abs(float64(weightMap[i] - expected)))
		fmt.Printf("backends[%d], %d, %d\n", i, weightMap[i], margin)
		if margin > maxMargin{
			maxMargin = margin
		}
	}

	fmt.Printf("expected:%d, maxMargin:%d\n", clientSize*subsetSize/backendSize, maxMargin)
}

func test4(){

	// 假定IP的范围是 192.168.1.1 ~ 192.168.1.255之间
	var loadMap map[int]int = make(map[int]int)

	for i:=0;i<clientSize;i++{
		ip := fmt.Sprintf("192.168.1.%d", i)
		// init backend
		backends:= make([]int, backendSize)
		for i:=0;i<backendSize;i++{
			backends[i] = i+1
		}
		//clientId := calcuCrc32(ip)
		//clientId := i + 53
		//clientId := calcuMurmur32(ip)
		//clientId := rand.Intn(100)

		clientId := int(InetAtoN(ip))
		fmt.Println("clientId:", clientId)

		setSetSlice := subset(backends, clientId, subsetSize)
		for _, backendId:=range setSetSlice{
			loadMap[backendId] = loadMap[backendId] + clientTotalReq/ len(setSetSlice)
		}
	}
	expected := clientTotalReq*clientSize/backendSize
	maxMargin := 0

	// print result
	for i:=1;i<backendSize+1;i++{
		margin := int(math.Abs(float64(loadMap[i] - expected)))
		fmt.Printf("backends[%d], %d, %d\n", i, loadMap[i], margin)
		if margin > maxMargin{
			maxMargin = margin
		}
	}

	fmt.Printf("expected:%d, maxMargin:%d\n", expected, maxMargin)
}

func main() {
	//test2()
	//test()
	//test3()
	test4()
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func calcuCrc32(str string) int{
	hasher := crc32.NewIEEE()
	hasher.Write([]byte(str))
	return int(hasher.Sum32())
}

func calcuMurmur32(str string) int{
	h32 := murmur3.New32()
	h32.Write([]byte(str))
	return int(h32.Sum32())
}

func calcuMurmur64(str string) int{
	h32 := murmur3.New64()
	h32.Write([]byte(str))
	x := int(h32.Sum64())
	if x < 0{
		return -x
	}
	return x
}


// 采用的方案是<<SRE Google运维解密>>书中提到的子集选择算法二：确定性子集。
// clientId是将原client的Ip地址做了CRC处理，转成int
// subsetSize就是子集的大小
func subset(backends []int, clientId int, subsetSize int) []int {
	subSetCount := len(backends) / subsetSize

	round := clientId / subSetCount
	rand.Seed(int64(round))
	rand.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})

	subsetId := clientId % subSetCount
	start := subsetId * subsetSize

	fmt.Println("round", round, "start", start)
	return backends[start : start+subsetSize]
}

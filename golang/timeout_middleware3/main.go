package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/vearne/gin"
	"github.com/vearne/golib/buffpool"
	"github.com/vearne/golib/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	HandlerFuncTimeout = "E501"

)

var greeter pb.GreeterClient

func init(){


	func init() {
		encoding.RegisterCodec(protoCodec{})
	}
	conn, err := grpc.Dial(address,
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	greeter = pb.NewGreeterClient(conn)
}

type errResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}


type TimeoutWriter struct {
	gin.ResponseWriter
	// body
	body *bytes.Buffer
	// header
	h http.Header

	mu sync.Mutex
	timedOut    bool
	wroteHeader bool
	code int
}

func (tw *TimeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		//return 0, http.ErrHandlerTimeout
		// 已经超时了，就不再写数据
		return 0, nil
	}

	return tw.body.Write(b)
}

func (tw *TimeoutWriter) WriteHeader(code int){
	fmt.Println("----xxx---", "TimeoutWriter-WriteHeader")
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *TimeoutWriter) WriteHeaderNow(){
	fmt.Println("----xxx---", "TimeoutWriter-WriteHeaderNow")
}

func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}


func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		// sync.Pool
		buffer := buffpool.GetBuff()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw


		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)

		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					fmt.Println("handler error", p, string(utils.Stack()))
					panicChan <- p
				}
			}()

			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
			bt, _ := json.Marshal(errResponse{Code: HandlerFuncTimeout,
				Msg: http.ErrHandlerTimeout.Error()})
			tw.ResponseWriter.Write(bt)
			c.Abort()
			cancel()
			tw.timedOut = true
			// 如果超时的话，buffer无法主动清除，只能等待GC回收
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			fmt.Println("tw.code", tw.code)
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(buffer.Bytes())
			buffpool.PutBuff(buffer)
		}
	}
}


func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func nocontent(c *gin.Context) {
	//c.Status(204)
	time.Sleep(1 * time.Second)
	c.Data(http.StatusNoContent, "", []byte{})
}

func long(c *gin.Context) {
	name := defaultName
	ctx := c.Request.Context()
	r, err := greeter.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Printf("could not greet: %v\n", err)
		return
	}
	log.Printf("Greeting: %s", r.Message)
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func main() {
	// create new gin without any middleware
	engine := gin.Default()

	// add timeout middleware with 2 second duration
	engine.Use(Timeout(time.Second * 1))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a route that will last 5 seconds
	engine.GET("/long", long)

	engine.GET("/nocontent", nocontent)

	// run the server
	log.Fatal(engine.Run(":8080"))
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vearne/golib/buffpool"
	"github.com/vearne/golib/utils"
	"google.golang.org/grpc"
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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

func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}


func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sync.Pool
		buffer := buffpool.GetBuff()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)

		finish := make(chan struct{})
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
			//tw.ResponseWriter.WriteHeader(401)
			tw.ResponseWriter.Write(buffer.Bytes())
			buffpool.PutBuff(buffer)
		}
	}
}


func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func long(c *gin.Context) {
	// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	////defer conn.Close()
	//greeter := pb.NewGreeterClient(conn)
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

	// run the server
	log.Fatal(engine.Run(":8080"))
}

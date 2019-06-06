package main

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vearne/golib/buffpool"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"net/http"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

type SimplebodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w SimplebodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sync.Pool
		buffer := buffpool.GetBuff()

		blw := &SimplebodyWriter{body: buffer, ResponseWriter: c.Writer}
		c.Writer = blw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)

		finish := make(chan struct{})
		go func() {
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			c.Abort()
			// 如果程序block在gRPC上，则终止执行
			cancel()
			// 如果超时的话，buffer无法主动清除，只能等待GC回收
		case <-finish:
			// 结果只会在主协程中被写入
			blw.ResponseWriter.Write(buffer.Bytes())
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	greeter := pb.NewGreeterClient(conn)
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
	engine := gin.New()

	// add timeout middleware with 2 second duration
	engine.Use(Timeout(time.Second * 1))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a route that will last 5 seconds
	engine.GET("/long", long)

	// run the server
	log.Fatal(engine.Run(":8080"))
}

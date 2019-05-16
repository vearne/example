package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/vearne/golib/buffpool"
	"log"
	"net/http"
	"time"
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

		finish := make(chan struct{})
		go func() {
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case <-time.After(t):
			c.Writer.WriteHeader(http.StatusGatewayTimeout)
			c.Abort()
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
	c.JSON(http.StatusOK, gin.H{"hello":"world"})
}

func long(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello":"world"})
}


func main() {
	// create new gin without any middleware
	engine := gin.New()

	// add timeout middleware with 2 second duration
	engine.Use(Timeout(time.Second * 2))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a route that will last 5 seconds
	engine.GET("/long", long)

	// run the server
	log.Fatal(engine.Run(":8080"))
}
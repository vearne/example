package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var jobChannel chan int

func init() {
	jobChannel = make(chan int, 10)
}

func main() {
	for i := 0; i < 3; i++ {
		jobChannel <- i
	}

	// 已包含register动作
	jobChanLength := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "job_channel_length",
		Help: "length of job channel",
	})
	go func() {
		jobChanLength.Set(float64(len(jobChannel)))
		c := time.Tick(30 * time.Second)
		for range c {
			jobChanLength.Set(float64(len(jobChannel)))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("starting...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

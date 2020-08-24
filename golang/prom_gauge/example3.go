package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var jobChannel chan int

func init() {
	jobChannel = make(chan int, 10)
}

func main() {
	for i := 0; i < 5; i++ {
		jobChannel <- i
	}

	jobChanLength := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "job_channel_length",
		Help: "length of job channel",
	},func() float64{
		return float64(len(jobChannel))
	})

	prometheus.MustRegister(jobChanLength)


	http.Handle("/metrics", promhttp.Handler())
	log.Println("starting...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

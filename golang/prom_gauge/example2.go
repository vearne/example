package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type ChanCollector struct{
	JobChannel chan int
	Name string
}

func NewChanCollector(jobChannel chan int, name string) *ChanCollector{
	c := ChanCollector{}
	c.JobChannel = jobChannel
	c.Name = name
	return &c
}

func (c *ChanCollector) Describe(ch chan<- *prometheus.Desc) {
	desc := prometheus.NewDesc(c.Name, "length of channel", nil, nil)
	ch <- desc
}

func (c *ChanCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(c.Name, "length of channel", nil, nil),
		prometheus.GaugeValue,
		float64(len(c.JobChannel)),
	)
}


var jobChannel chan int

func init() {
	jobChannel = make(chan int, 10)
}

func main() {
	for i := 0; i < 4; i++ {
		jobChannel <- i
	}

	collector := NewChanCollector(jobChannel, "job_channel_length")
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("starting...")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

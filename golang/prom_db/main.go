package main

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vearne/golib/metric"
	"log"
	"net/http"
	"time"
)


func main() {
	// init redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 100,
	})

	metric.AddRedis(client, "car")

	// init mysql
	DSN := "test:xxxx@tcp(localhost:6379)/somebiz?charset=utf8&loc=Asia%2FShanghai&parseTime=true"
	mysqldb, err := gorm.Open("mysql", DSN)
	if err != nil {
		panic(err)
	}

	mysqldb.DB().SetMaxIdleConns(50)
	mysqldb.DB().SetMaxOpenConns(100)
	mysqldb.DB().SetConnMaxLifetime(5 * time.Minute)
	//mysqldb = mysqldb.Debug()

	metric.AddMySQL(mysqldb, "car")

	// do some thing
	for i := 0; i < 30; i++ {
		go func() {
			for {
				client.Get("a").String()
				time.Sleep(200 * time.Millisecond)
				mysqldb.Exec("show tables")
			}
		}()
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9090", nil))
	log.Println("starting...")
}

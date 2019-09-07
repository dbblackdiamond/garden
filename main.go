package main

import (
  "log"
  "net/http"
  "github.com/influxdata/influxdb1-client/v2"
  "time"
  "github.com/garyburd/redigo/redis"
)

var (
  influxDBConn client.Client
  Pool *redis.Pool
  apiToken = "BBFF-8tET4fVLoB5RjFrpsQe7sCImQyKBsg"
  garden Garden
)

func check (function string, err error) {
  if err != nil {
    log.Fatalln(function, err)
  }
}

func main() {
  log.Println("Initializing Redis connection")
  redisInit()
  log.Println("Loading Garden config")
  loadGarden()
  log.Println("Starting garden server @ ", time.Now())
  router := NewRouter()

  log.Fatal(http.ListenAndServe(":8082", router))
}

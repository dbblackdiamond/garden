package main

import (
  "log"
  "net/http"
  "github.com/influxdata/influxdb1-client/v2"
  "time"
)

var (
  influxDBConn client.Client
)

func main() {
  log.Println("Starting garden server @ ", time.Now())
  router := NewRouter()

  log.Fatal(http.ListenAndServe(":8082", router))
}

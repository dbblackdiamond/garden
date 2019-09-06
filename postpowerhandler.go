package main

import (
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

func postPowerHandler(w http.ResponseWriter, r*http.Request) {
  var power PowerMsg
  var record Record

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&power)
  if err != nil {
      panic(err)
  }
  fmt.Println("\nReceived Power message from box ", power.Box, " with voltage ", power.Voltage)
  record.Box = power.Box
  record.Sensor = 0
  record.Value = power.Voltage
  record.Timestamp = int64(power.Timestamp)
  influxDBConn, err := initializeDBConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  err = insertRecords(influxDBConn, record, "power")
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  influxDBConn.Close()
}

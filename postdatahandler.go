package main

import (
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

/*
{"timestamp": uint32, "box": int, "sensor": int, "value": float}
*/

func postDataHandler(w http.ResponseWriter, r *http.Request) {
  var record Record
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&record)
  if err != nil {
      panic(err)
  }
  record.VWC = convertVoltageToVWCFloat(record.Value)
  fmt.Println("postDataHandler>\tReceived record: timestamp = ", record.Timestamp, ", box = ", record.Box, ", sensor = ", record.Sensor, ", value = ", record.Value, " VWC = ", record.VWC)

  influxDBConn, err := initializeDBConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  err = insertRecords(influxDBConn, record, "moisture")
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  influxDBConn.Close()
}

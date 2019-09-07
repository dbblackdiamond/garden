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

  db, err := initializePQConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  sqlStatement := `INSERT INTO moisture(box, sensor, value, vwc) VALUES($1, $2, $3, $4)`
  err = db.Exec(sqlStatement, record.Box, record.Sensor, record.Value, record.VWC)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  db.Close()
}

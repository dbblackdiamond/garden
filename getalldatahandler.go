package main

import (
  "encoding/json"
  "net/http"
  "log"
  "io"
)

func getAllDataHandler(w http.ResponseWriter, r *http.Request) {
  influxDBConn, err := initializeDBConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  resQuery, err := getAllRecords(influxDBConn)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  results := analyzeRecords(influxDBConn, resQuery)

  JSONResults, err := json.Marshal(results)
  if err != nil {
    log.Fatalln("Error: ", err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, string(JSONResults))
  influxDBConn.Close()
}

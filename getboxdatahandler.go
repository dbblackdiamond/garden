package main

import (
  "encoding/json"
  "net/http"
  "io"
  "github.com/gorilla/mux"
  "log"
  "strconv"
)

func getBoxDataHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  box, _ := strconv.Atoi(vars["box"])

  influxDBConn, err := initializeDBConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  resQuery, err := getBoxRecords(influxDBConn, box)
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

package main

import (
  "net/http"
  "io"
  "log"
  "fmt"
  "strconv"
  "github.com/gorilla/mux"
)

func getWateringDecision(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  box, _ := strconv.Atoi(vars["box"])

  conn, err := initializeDBConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  decision, err := getWateringRecord(conn, box)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  conn.Close()

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  if decision == 1 {
    fmt.Println("getWateringDecision: decision = ", decision, " so watering is true")
    io.WriteString(w, `{"watering": true}`)
  } else {
    fmt.Println("getWateringDecision: decision = ", decision, " so watering is false")
    io.WriteString(w, `{"watering": false}`)
  }
}

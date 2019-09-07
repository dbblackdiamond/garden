package main

import (
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

func postPowerHandler(w http.ResponseWriter, r*http.Request) {
  var power PowerMsg

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&power)
  if err != nil {
      panic(err)
  }
  fmt.Println("\nReceived Power message from box ", power.Box, " with voltage ", power.Voltage)
  db, err := initializePQConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  sqlStatement := `INSERT INTO power(box, capacity, charge, current, health, voltage, signal) VALUES($1, $2, $3, $4, $5, $6, $7)`
  err = db.Exec(sqlStatement, power.Box, power.Capacity, power.Charge, power.Current, power.Health, power.Voltage, power.Signal)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  db.Close()
}

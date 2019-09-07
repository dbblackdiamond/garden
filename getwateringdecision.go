package main

import (
  "net/http"
  "io"
  "log"
  "strconv"
  "github.com/gorilla/mux"
  "github.com/garyburd/redigo/redis"
  "fmt"
)

func getWateringDecision(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := ""
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      key = garden.Devices[i].Key
    }
  }
  if key == "" {
    log.Fatalln("getWateringDecision:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  wateringStr, err := redis.Strings(c.Do("HMGET", key, "watering"))
  check("Hmget", err)
  decision, _ := strconv.Atoi(wateringStr[0])

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

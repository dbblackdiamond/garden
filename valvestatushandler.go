package main

import (
  "log"
  "net/http"
  "io"
  "github.com/gorilla/mux"
  "github.com/garyburd/redigo/redis"
)

//Improvement to do: make it persistent with either influxdb or redis!!!!!
func valveStatusHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := ""
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      key = garden.Devices[i].Key
    }
  }
  if key == "" {
    log.Fatalln("valveStatusHandler:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  _, err := c.Do("HMSET", key, "valve", vars["value"])
  check("Hmset", err)
}

func getValveStatusHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := ""
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      key = garden.Devices[i].Key
    }
  }
  if key == "" {
    log.Fatalln("getValveStatusHandler:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  valve, err := redis.Strings(c.Do("HMGET", key, "valve"))
  check("Hmget", err)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  jsonStr := `{"valve":` + valve[0] + `}`
  io.WriteString(w, jsonStr)
}

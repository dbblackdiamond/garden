package main

import (
  "log"
  "net/http"
  "io"
  "github.com/gorilla/mux"
  "github.com/garyburd/redigo/redis"
)

//Improvement to do: make it persistent with either influxdb or redis!!!!!
func getBoxNumberHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := ""
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      key = garden.Devices[i].Key
    }
  }
  if key == "" {
    log.Fatalln("getBoxNumberHandler:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  box, err := redis.Strings(c.Do("HMGET", key, "box"))
  check("Hmget", err)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  jsonStr := `{"box":` + box[0] + `}`
  io.WriteString(w, jsonStr)
}

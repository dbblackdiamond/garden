package main

import (
  "net/http"
  "io"
  "log"
  "strconv"
  "github.com/gorilla/mux"
  "github.com/garyburd/redigo/redis"
  "fmt"
  "encoding/json"
)

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  hash := ""
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      hash = garden.Devices[i].Key
    }
  }
  if hash == "" {
    log.Fatalln("getConfigHandler:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  keys, err := redis.Strings(c.Do("HKEYS", hash))
  check("Hkeys", err)
  jsonMap  = map[string]interface{}
  for _, key := range keys {
    jsonMap[key], err := redis.Strings(c.Do("HMGET", hash, key))
    check("Hmget", err)
  }

  jsonStr, _ := json.Marshall(jsonMap)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  fmt.Println("getconfigHandler: config = ", jsonMap)
  io.WriteString(w, jsonStr)
}

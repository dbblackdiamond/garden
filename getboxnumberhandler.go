package main

import (
  "net/http"
  "io"
  "github.com/gorilla/mux"
  "strconv"
  "fmt"
)

var boxes map[int]string

//Improvement to do: make it persistent with either influxdb or redis!!!!!
func getBoxNumberHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  for key, value := range boxes {
    fmt.Println("getBoxNumber: key = ", key, " and value = ", value)
    if(value == vars["deviceid"]) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(http.StatusOK)
      fmt.Println("getBoxNumber: for deviceid ", value, " the box number is ", key)
      jsonStr := `{"box":` + strconv.Itoa(key) + `}`
      io.WriteString(w, jsonStr)
      return
    }
  }
  fmt.Println("getBoxNumber: deviceid not found, adding it")
  index := len(boxes)
  if index == 0 {
    boxes = make(map[int]string)
  }
  //Need to avoid box 0 as Influxdb client doesn't handle it properly.
  boxes[index + 1] = vars["deviceid"]
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  jsonStr := `{"box":` + strconv.Itoa(index + 1) + `}`
  io.WriteString(w, jsonStr)

}

func addBoxNumberHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  box, _ := strconv.Atoi(vars["box"])
  if boxes == nil {
    boxes = make(map[int]string)
  }
  boxes[box] = vars["deviceid"]
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  fmt.Println("addBoxNumber: for deviceid ", boxes[box], " the box number is ", box)
  jsonStr := `{"box":` + vars["box"] + `}`
  io.WriteString(w, jsonStr)
}

func deleteBoxesHandler(w http.ResponseWriter, r *http.Request) {
  boxes = nil
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  jsonStr := `{"boxdeleted": true}`
  io.WriteString(w, jsonStr)
}

package main

import (
  "net/http"
  "github.com/gorilla/mux"
)

type Route struct {
  Name string
  Method string
  Pattern string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    var handler http.Handler
    handler = route.HandlerFunc
    handler = Logger(handler, route.Name)

    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(handler)
  }
  return router
}

var routes = Routes{
  Route{ "Index", "GET", "/", indexHandler},
  Route{ "PostData", "POST", "/api/postdata", postDataHandler},
  Route{ "PostPower", "POST", "/api/power", postPowerHandler},
  Route{ "GetAllData", "GET", "/api/getdata", getAllDataHandler},
  Route{ "GetBoxData", "GET", "/api/getdata/{box}", getBoxDataHandler},
  Route{ "Watering", "GET", "/api/watering/{box}", getWateringDecision},
  Route{ "GetBoxNumber", "GET", "/api/getboxnumber/{deviceid}", getBoxNumberHandler},
  Route{ "AddBoxNumber", "POST", "/api/addbox/{box}/{deviceid}", addBoxNumberHandler},
  Route{ "DeleteBoxes", "POST", "/api/deleteboxes", deleteBoxesHandler},
}

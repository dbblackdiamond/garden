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

func getDoAnalysisHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := ""
  devicepos := 0
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == vars["deviceid"] {
      key = garden.Devices[i].Key
      devicepos = i
    }
  }
  if key == "" {
    log.Fatalln("getDoAnalysisHandler:\tCan't find device ", vars["deviceid"])
  }

  c := Pool.Get()
  defer c.Close()

  box, err := redis.Strings(c.Do("HMGET", key, "box"))
  check("Hmget", err)
  box := strconv.Itoa(box[0])

  fmt.Println("main: Getting weather data")
  weatherData := getWeatherData()
  fmt.Println("main: Analysing weather data")
  weatherDecision := weatherDataAnalysis(weatherData)
  if garden.Devices[devicepos].Valve == 1 {
    checkWateringDuration(box, deviceid, devicepos)
    //if we are already watering the garden box
    //get the time at which watering was started and if more than 10 mins, cut it off.
  } else {
    //if we are not watering the box, then we need to figure out if we need to
    fmt.Println("main: Getting sensor data for box", box)
    if box != 0 {
      sensors := getSensorData(deviceid, box)
      counter := 0
      for _, sensor := range sensors {
        fmt.Println("main: counter =", counter, " sensor.Water = ", sensor.Water)
        counter += sensor.Water
      }

      decision := 0
      avg := int(counter/garden.Devices[devicepos].NumSensors)
      if force == 1 {
        avg = 4
      }
      switch avg {
        //Garden is very wet, no water
        case 0:
          fmt.Println("main: avg = ", avg, " and wet, no watering needed")
        //Garden is optimal, but slightly wet, no water
        case 1:
          fmt.Println("main: avg = ", avg, " slightly wet, no watering needed")
        //Garden is optimal, but slightly dry, watering depends on weather
        case 2:
          //No rain in the forecast, we water
          if weatherDecision == false {
            fmt.Println("main: avg = ", avg, " and no rain in the forecast, watering needed")
            decision = 1
          } else { //Rain in the forecast, no need to water
            fmt.Println("main: avg = ", avg, " and rain in the forecast, no watering needed")
          }
        //Garden is slightly dry, but no too much, watering depends on weather
        case 3:
          //No rain in the forecast, we water
          if weatherDecision == false {
            fmt.Println("main: avg = ", avg, " and no rain in the forecast, watering needed")
            decision = 1
          } else { //Rain in the forecast, no need to water
            fmt.Println("main: avg = ", avg, " and rain in the forecast, no watering needed")
          }
        //Garden is very dry, watering mandatory
        case 4:
          fmt.Println("main: avg = ", avg, " and very dry, watering mandatory")
          decision = 1
      }
      c := Pool.Get()
      defer c.Close()
      key := getKey(deviceid)
      _, err := c.Do("HMSET", key, "watering", decision )
      check("hmset", err)
      garden.Devices[devicepos].StartedAt = time.Now().Unix()
      garden.Devices[devicepos].Watering = decision
    }
  }
}

func checkWateringDuration(box int, deviceid string, devicepos int) {
  c := Pool.Get()
  defer c.Close()

  fmt.Println("main: Box ", box, " is already being watered, checking duration")
  db, err := initializePQConnection()
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  defer db.Close()
  sqlStatement := `select created_at from valve `
  err = db.Exec(sqlStatement, record.Box, record.Sensor, record.Value, record.VWC)
  if err != nil {
    log.Fatalln("Error: ", err)
  }

  timeDiff := time.Now().Unix() - garden.Devices[devicepos].StartedAt
  //if we have been watering for more than 5 minutes, we cut it off.
  //Otherwise, we do nothing and let it run
  if timeDiff > wateringDuration {
    fmt.Println("main: Been watering for more than ", wateringDuration, ", stopping it")
    c := Pool.Get()
    defer c.Close()
    key := getKey(deviceid)
    _, err := c.Do("HMSET", key, "valve", 0 )
    check("hmset", err)
    garden.Devices[devicepos].Valve = 0
  } else {
    fmt.Println("main: Watering duration not expired, watering continues")
  }
}

package main

import (
//  "fmt"
  "strconv"
  "log"
  "github.com/influxdata/influxdb1-client/v2"
)

func findDeviceid(box int) (string, int) {
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Box == box {
      //fmt.Println("main: device id = ", garden.Devices[i].Id, " for box ", box)
      return garden.Devices[i].Id, i
    }
  }
  return "", 0
}

func getKey(deviceid string) string {
  for i := 0; i < garden.NumDevices; i++ {
    if garden.Devices[i].Id == deviceid {
      //fmt.Println("main: key = ", garden.Devices[i].Key, " for device id ", deviceid)
      return garden.Devices[i].Key
    }
  }
  return ""
}

func convertVoltageToVWC(samples Samples) []float64 {
  vwc := make([]float64, len(samples))
  for i, sample := range samples {
    switch {
    case sample.Value >= 0.0 && sample.Value <= 1.1:
      vwc[i] = (sample.Value * 10.0) - 1.0
    case sample.Value > 1.1 && sample.Value <= 1.3:
      vwc[i] = (sample.Value * 25.0) - 17.5
    case sample.Value > 1.3 && sample.Value <= 1.82:
      vwc[i] = (sample.Value * 48.08) - 47.5
    case sample.Value > 1.82 && sample.Value <= 2.2:
      vwc[i] = (sample.Value * 26.32) - 7.89
    case sample.Value > 2.2 && sample.Value <= 3.0:
      vwc[i] = (sample.Value * 62.5) - 87.5
    default:
      vwc[i] = -1.0
    }
  }
  //fmt.Println("convertVoltageToVWC: after conversion, the VWC array is: ", vwc)
  return vwc
}

func calculateVWCAverage(data []float64) float64 {
  avg := 0.0
  for i := 0; i < len(data); i++ {
//    fmt.Println("calculateVWCAverage: avg = ", avg, " and data[i] = ", data[i])
    avg += data[i]
  }
  avg /= float64(len(data))
  //fmt.Println("calculateVWCAverage: the average VWC is: ", avg)
  return avg
}

//returns 0 if moisture value is above 51%, ie garden is too wet
//returns 1 if moisture value is between 45% and 51%, ie garden is good, no need to watering
//returns 2 if moisture value is between 40% and 45%, ie garden is good, but
//  if weather is dry, then we water
//  if weather calls for rain, then no water
//returns 3 if moisture value is between 30% and 40%, ie garden is dry, but not too dry, so
//  if weather is dry, then we water
//  if weather calls for rain, then no water
//returns 4 if moisture value is below 30%, ie very dry and this will force watering now matter the weather
func analyzeSensorData(sensor Sensor) int {
  vwc := convertVoltageToVWC(sensor.Samples)
  avg := calculateVWCAverage(vwc)
  switch {
  case avg >= 0.0 && avg < 30.0:
//    fmt.Println("detectSensorTrend: sensor #", sensor.Number, " average VWC is below 30%, returning 4")
    return 4
  case avg <= 30.0 && avg < 40.0:
//    fmt.Println("detectSensorTrend: sensor #", sensor.Number, " average VWC is between 30% and 40%, return 3")
    return 3
  case avg <= 40.0 && avg < 45.0:
//    fmt.Println("detectSensorTrend: sensor #", sensor.Number, " average VWC is between 40% and 45%, return 3")
    return 2
  case avg >= 45.0 && avg < 51.0:
//    fmt.Println("detectSensorTrend: #", sensor.Number, " average VWC is between 40% and 51%, returning 1")
    return 1
  case avg >= 51.0:
//    fmt.Println("detectSensorTrend: #", sensor.Number, " average VWC is above 51%, returning 0")
    return 0
  }
  return 0
}

func getSensorData(deviceid string, box int) Sensors {
  var sensors Sensors
  for i := 0; i < garden.Devices[devicepos].NumSensors; i++ {
    var sensor Sensor
    results := getSensorRecord(deviceid, i)
    sensor.Number = i
    sensor.Dots = results
    sensors = append(sensors, sensor)
  }
  for i := 0; i < garden.Devices[devicepos].NumSensors; i++ {
    res := analyzeSensorData(sensors[i])
    //fmt.Println("getSensorData: pre-switch sensor.Water = ", sensors[i].Water)
    switch res {
    case 0:
      sensors[i].Water = 0
    case 1:
      sensors[i].Water = 1
    case 2:
      sensors[i].Water = 2
    case 3:
      sensors[i].Water = 3
    case 4:
      sensors[i].Water = 4
    }
    //fmt.Println("getSensorData: pre-switch sensor.Water = ", sensors[i].Water)
  }
  //fmt.Printf("getSensorData: sensors = %+v\n", sensors)
  return sensors
}

package main

import (
//  "fmt"
  "strconv"
  "log"
  "github.com/influxdata/influxdb1-client/v2"
)

type Sample struct {
  Timestamp int64 `json:"t"`
  Value float64 `json:"v"`
}

type Samples []Sample

type Sensor struct {
  Number int `json:"slot"`
  Samples []Sample `json:"samples"`
  Water int `json:"water"`
}

type Sensors []Sensor

func getNumberOfBoxes(conn client.Client) []int {
  results, err := getListOfBoxes(conn)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  boxes := make([]int, len(results[0].Series[0].Values))
  for i, value := range results[0].Series[0].Values {
    //fmt.Println("getNUmberOfBoxes: Record #",i ,":\tBox = ", value[1])
    val, _ := value[1].(string)
    boxes[i], _ = strconv.Atoi(val)
  }
  //fmt.Println("getNUmberOfBoxes: Found ", len(boxes), " boxes to analyze")
  return boxes
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

func getNumberOfSensors(records Records) int {
  counter := make(map[int]int)
  for _, row := range records {
      counter[row.Sensor]++
    }
    //fmt.Println("getNumberOfSensors: Found ", len(counter), " sensors in the records")
    return len(counter)
}

func convertRecordIntoSensorData(records Records) Sensors {
  //numberSensors := getNumberOfSensors(records)
  sensors := make(Sensors, 1)
  var sample Sample
  for i := 0; i < len(records); i++ {
    //fmt.Println("convertRecordIntoSensorData: Converting record #", i, " out of ", len(records), " records")
    //fmt.Println("convertRecordIntoSensorData: len(sensor) = ", len(sensors), " and records[i].Sensor = ", records[i].Sensor)
    if len(sensors) > records[i].Sensor {
      //fmt.Println("convertRecordIntoSensorData: appending to existing sensor data")
      sample.Value = records[i].Value
      sample.Timestamp = records[i].Timestamp
      sensors[records[i].Sensor].Samples = append(sensors[records[i].Sensor].Samples, sample)
    } else {
      var sensor Sensor
      //fmt.Println("convertRecordIntoSensorData: creating new sensor")
      sensor.Number = records[i].Sensor
      sensor.Samples = make(Samples, 1)
      sensor.Samples[0].Value = records[i].Value
      sensor.Samples[0].Timestamp = records[i].Timestamp
      sensors = append(sensors, sensor)
    }
  }
  //fmt.Println("ConvertRecordIntoSensorData: Converted records ", records, " into sensors ", sensors)
  return sensors
}

func getSensorData(conn client.Client, box int) Sensors {
  res, err := getBoxRecords(conn, box)
  if err != nil {
    log.Fatalln("Error: ", err)
  }
  //fmt.Println("getSensorData: Analyzing records")
  records := analyzeRecords(conn, res)
  //fmt.Println("getSensorData: Converting records into sensor data")
  sensors := convertRecordIntoSensorData(records)
  for i := 0; i < len(sensors); i++ {
    sensors[i].Water = analyzeSensorData(sensors[i])
  }
  //fmt.Printf("getSensorData: sensors = %+v\n", sensors)
  return sensors
}

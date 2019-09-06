package main

import (
  "github.com/influxdata/influxdb1-client/v2"
  "encoding/json"
  "strconv"
  "fmt"
  "time"
)

func getWateringRecord(conn client.Client, box int) (int64, error) {
  q := client.Query{
    Command: fmt.Sprintf("select * from watering where box = '%d' order by desc limit 1", box),
    Database: database,
  }

  resp, err := conn.Query(q)
  if err != nil {
    return -1, err
  }
  if resp.Error() != nil {
    return -1, resp.Error()
  }
  if len(resp.Results[0].Series[0].Values) != 1 {
    return -1, nil
  }
  tempTime, _ := time.Parse(time.RFC3339, resp.Results[0].Series[0].Values[0][0].(string))
  recordTime := int64(tempTime.Unix())
  currentTime := time.Now().UTC().Unix()
  if currentTime - recordTime <= 7200 {
    fmt.Println("getWateringRecord: record is less than 2 hours old, getting decision")
    decision, _ := resp.Results[0].Series[0].Values[0][2].(json.Number).Int64()
    return decision, nil
  }
  fmt.Println("getWateringRecord: record is more than 2 hours old, no decision")
  return -1, nil
}

func analyzeRecords(c client.Client, results []client.Result) []Record {
  records := make([]Record, len(results[0].Series[0].Values))
  for i, value := range results[0].Series[0].Values {
//    fmt.Println("Record #",i ,":\tBox = ", value[1], "\tSensor = ", value[2], "\tValue = ", value[3], "\tTimestamp = ", value[0])
    timestamp, _ := time.Parse(time.RFC3339, value[0].(string))
    records[i].Timestamp = int64(timestamp.Unix())
    box, _ := value[1].(string)
    records[i].Box, _ = strconv.Atoi(box)
    sensor, _ := value[2].(string)
    records[i].Sensor, _ = strconv.Atoi(sensor)
    records[i].Value, _ = value[3].(json.Number).Float64()
  }
  return records
}

func convertVoltageToVWCFloat(voltage float64) float64 {
  switch {
  case voltage >= 0.0 && voltage <= 1.1:
    return ((voltage * 10.0) - 1.0)
  case voltage >= 1.1 && voltage <= 1.3:
    return ((voltage * 25.0) - 17.5)
  case voltage >= 1.3 && voltage <= 1.82:
    return ((voltage * 48.08) - 47.5)
  case voltage >= 1.82 && voltage <= 2.2:
    return ((voltage * 26.32) - 7.89)
  case voltage >= 2.2 && voltage <= 3.0:
    return ((voltage * 62.5) - 87.5)
  default:
    return -1.0
  }
  return -1.0
}

func insertRecords(c client.Client, record Record, measurement string) error {
  var tags map[string]string
  bp, err := client.NewBatchPoints(client.BatchPointsConfig{
    Database: database,
    Precision: "s",
  })
  if err != nil {
    return err
  }
  if measurement == "moisture" {
    tags = map[string]string{
      "box": strconv.Itoa(record.Box),
      "sensor": strconv.Itoa(record.Sensor),
      "vwc": fmt.Sprintf("%.2f", record.VWC),
    }
  } else {
    tags = map[string]string{
      "box": strconv.Itoa(record.Box),
      "sensor": strconv.Itoa(record.Sensor),
    }
  }
  fields := map[string]interface{} {
    "value": record.Value,
  }
  point, err := client.NewPoint(
    measurement,
    tags,
    fields,
    time.Unix(record.Timestamp, 0),
  )
  if err != nil {
    return err
  }
  bp.AddPoint(point)

  err = c.Write(bp)
  if err != nil {
    return err
  }
  return nil
}

func getAllRecords(c client.Client) ([]client.Result, error) {
  q := client.Query{
    Command: "select * from moisture",
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}

func getBoxRecords(c client.Client, box int) (res []client.Result, err error) {
  q := client.Query{
    Command: fmt.Sprintf("select * from moisture where box = '%d' AND time > now() - 24h tz('America/Edmonton')", box),
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}

func getListOfBoxes(c client.Client) ([]client.Result, error) {
  q := client.Query{
    Command: "select distinct(box) from (select box, value from moisture)",
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}

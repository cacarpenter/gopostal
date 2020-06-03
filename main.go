package main

import (
  "encoding/json"
  "fmt"
  "github.com/cacarpenter/gopostal/gp"
  "io/ioutil"
)

const TESTFILE = "/home/ccarpenter/Documents/postman/carrier_postman_collection.json"

func main() {
  data, err := ioutil.ReadFile(TESTFILE)
  if err != nil {
    panic(err)
  }

  var coll gp.PostmanCollection
  err = json.Unmarshal(data, &coll)
  if err != nil {
    panic(err)
  }

  fmt.Println(coll)
}

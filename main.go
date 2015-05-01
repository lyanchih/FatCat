package main

import (
  "./mega"
  "fmt"
)

func main() {
  url := "https://mega.co.nz/#!sQ4HBA6a!md9EZHtkl_-A9hSJIy2KR4_4tFAZ_1p5dxhE2_KBcjE"
  m := mega.New(url, "tmp")
  m.Download()

  if m.Error() != nil {
    fmt.Println(m)
  }
}

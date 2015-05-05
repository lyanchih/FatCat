package main

import (
  "./mega"
  "./task"
  "os"
  "fmt"
  "log"
  "io/ioutil"
)

const MaxWorker = 10

// var sample []mega.MegaInfo = []mega.MegaInfo{
//   {"https://mega.co.nz/#!sQ4HBA6a!md9EZHtkl_-A9hSJIy2KR4_4tFAZ_1p5dxhE2_KBcjE", "", 1},
//   {"https://mega.co.nz/#!sQ4HBA6a!md9EZHtkl_-A9hSJIy2KR4_4tFAZ_1p5dxhE2_KBcjE", "", 2},
// }

func main() {
  if len(os.Args) == 0 {
    log.Fatal("Please offer json file.")
  }
  pool, err := task.CreatePool()
  if err != nil {
    log.Fatal(err)
  }

  pool.Registry(task.MegaType, mega.MegaRoot{})

  file, err := os.Open(os.Args[1])
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  
  data, err := ioutil.ReadAll(file)
  if err != nil {
    log.Fatal(err)
  }
  
  _, err = pool.Add(task.MegaType, file.Name(), data)
  if err != nil {
    log.Fatal(err)
  }
  
  go pool.Start()

  for i := 0; i < MaxWorker; i++{
    go func(i int) {
      fmt.Println("Process", i, "asking...")
      t, ok := pool.Ask()
      if ok {
        fmt.Println("Process", i, "downloading...")
        t.Download()
        fmt.Println("Process", i, "reporting...")
        pool.Report(t)
        fmt.Println("Process", i, "done")
      }
    }(i)
  }

  select{}
}

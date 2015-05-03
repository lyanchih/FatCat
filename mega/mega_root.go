package mega

import (
  "../task"
  "fmt"
  "sort"
  "encoding/json"
)

func (m MegaRoot) Parse(data []byte) (ms []task.Tasker, err error) {
  var arr megaInfoArr
  err = json.Unmarshal(data, &arr)
  if err != nil {
    return nil, err
  } else if len(arr) == 0 {
    return nil, fmt.Errorf("Data didn't contain any mega information.")
  }
  
  sort.Sort(arr)
  ms = make([]task.Tasker, 0, len(arr))
  for i, m := range arr {
    if m.Part == 0 {
      m.Part = uint(i) + 1
    }
    if m.Name == "" {
      m.Name = fmt.Sprintf("part%d", i+1)
    }
    ms = append(ms, NewMega(m.Url, m.Name))
  }
  return ms, nil
}

type megaInfoArr []MegaInfo

func (arr megaInfoArr) Len() int {
  return len(arr)
}

func (arr megaInfoArr) Swap(i, j int) {
  arr[i].Url, arr[i].Name, arr[j].Url, arr[j].Name = arr[j].Url, arr[j].Name, arr[i].Url, arr[i].Name
}

func (arr megaInfoArr) Less(i, j int) bool {
  return arr[i].Part > arr[j].Part
}

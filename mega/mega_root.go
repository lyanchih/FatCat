package mega

import (
  "../task"
  "fmt"
  "encoding/json"
)

func (m MegaRoot) Parse(data []byte) (ms []task.Tasker, err error) {
  var arr []MegaInfo
  err = json.Unmarshal(data, &arr)
  if err != nil {
    return nil, err
  } else if len(arr) == 0 {
    return nil, fmt.Errorf("Data didn't contain any mega information.")
  }
  
  ms = make([]task.Tasker, 0, len(arr))
  for _, m := range arr {
    ms = append(ms, &Mega{url:m.Url})
  }
  return ms, nil
}

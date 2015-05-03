package task

import (
  "log"
)

func (tg *TaskGroup) IsDone() bool {
  return len(tg.preTask) == 0 && len(tg.postTask) == 0 && len(tg.pending) == 0
}

func (tg *TaskGroup) report(t Tasker) {
  for i, _t := range tg.pending {
    if _t != t {
      continue
    }
    
    if t.Error() != nil {
      log.Println(t.Error())
      tg.preTask, tg.pending = append(tg.preTask, t), append(tg.pending[0:i], tg.pending[i+1:]...)
    } else {
      log.Printf("Task name %s was finishing.", t.Name())
      tg.done, tg.pending = append(tg.done, t), append(tg.pending[0:i], tg.pending[i+1:]...)
    } 
  }
}

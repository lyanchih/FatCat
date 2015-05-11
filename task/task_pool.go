package task

import (
  "fmt"
  "time"
)

const MaxTaskQueue = 10

func CreatePool() (*TaskPool, error) {
  return &TaskPool{
    make(map[uint]*TaskGroup, 1),
    make(map[TaskGroupType]Parser, 1),
    make(chan Task, MaxTaskQueue),
    make(chan Task, MaxTaskQueue),
    false,
  }, nil
}

func (tp *TaskPool) Registry(t TaskGroupType, p Parser) {
  tp.parsers[t] = p
}

func (tp *TaskPool) Unregistry(t TaskGroupType) {
  tp.parsers[t] = nil
}

func (tp *TaskPool) Add(t TaskGroupType, name string, data []byte) (uint, error) {
  p, ok := tp.parsers[t]
  if !ok {
    return 0, fmt.Errorf("This task group type (%d) didn't been registry yet.", t)
  }

  taskers, err := p.Parse(data)
  if err != nil {
    return 0, err
  }

  id := uint(len(tp.groups))
  tp.groups[id] = &TaskGroup{t, name, taskers, make([]Tasker, 0, 1), make([]Tasker, 0, 1), make([]Tasker, 0, 1)}
  return id, nil
}

func (tp *TaskPool) Start() {
  if tp.isStart {
    return
  }

  defer func() {
    recover()
    tp.isStart = false
  }()

  var t Task
  for {
    if t.Tasker == nil {
      select {
      case reportT := <- tp.reportChannel:
        tp.reportTask(reportT)
      case <-time.Tick(5 * time.Second):
        t = tp.askTask()
      } 
    } else {
      select {
      case reportT := <- tp.reportChannel:
        tp.reportTask(reportT)
      case tp.askChannel <- t:
        t = tp.askTask()
      }
    }
  }
}

func (tp *TaskPool) Stop() {
  close(tp.askChannel)
  close(tp.reportChannel)
}

func (tp *TaskPool) Ask() (Task, bool) {
  t, ok := <- tp.askChannel
  return t, ok
}

func (tp *TaskPool) Report(t Task) {
  if tp.reportChannel == nil {
    return
  }
  
  defer func() {
    recover()
  }()
  tp.reportChannel <- t
}

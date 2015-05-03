package task

type TaskGroupType uint

const (
  MegaType TaskGroupType = iota
  BitType
)

type TaskPool struct {
  groups map[uint]*TaskGroup
  parsers map[TaskGroupType]Parser
  askChannel chan Task
  reportChannel chan Task
}

type TaskGroup struct {
  TaskGroupType
  Name string
  preTask, postTask, pending, done []Tasker
}

type Task struct {
  id uint
  Tasker
}

type Tasker interface {
  Download()
  Name() string
  Error() error
}

type Parser interface {
  Parse(bs []byte) ([]Tasker, error)
}

package task

func (tp *TaskPool) askTask() (_ Task) {
  for id, g := range tp.groups {
    if len(g.preTask) != 0 {
      t := g.preTask[0]
      g.preTask, g.pending = g.preTask[1:], append(g.pending, t)
      return Task{id, t}
    }
  }

  return Task{0, nil}
}

func (tp *TaskPool) reportTask(t Task) {
  for id, g := range tp.groups {
    if id != t.id {
      continue
    }

    g.report(t.Tasker)
  }
}

package main

import (
    "strconv"
    "sync/atomic"
    "sync"
    "fmt"
)

type Recipe []map[string]int

type Brewer struct {
    is_brewing int32
    task_id int
    recipe Recipe
    abort chan struct{}
    abort_mutex sync.Mutex
}

func (brewer *Brewer) Lock() bool {
    return atomic.CompareAndSwapInt32(&brewer.is_brewing, 0, 1)
}

func (brewer *Brewer) Unlock() {
    brewer.is_brewing = 0
}

func (brewer *Brewer) Brew(recipe Recipe) bool {
    if !brewer.Lock() { return false }
    brewer.recipe = recipe
    brewer.task_id++

    go brewer.BrewRoutine(recipe)

    return true;
}

func (brewer *Brewer) GetTaskId() int {
  if brewer.is_brewing == 1 {
    return brewer.task_id
  } else {
    return -1
  }
}

func (brewer *Brewer) BrewRoutine(recipe Recipe) {
    defer brewer.Unlock()

    brewer.abort = make(chan struct{})

    stepLoop:
    for stepi,step := range recipe {
        stepWg := sync.WaitGroup{}
        stepWg.Add(len(step))
        for idStr,time := range step {
          var id, _ = strconv.ParseInt(idStr, 10, 64)
          go ServePump(id, time, &stepWg, brewer.abort)
        }

        stepWg.Wait()
        select {
          case <- brewer.abort:
            break stepLoop
          default:
        }
        fmt.Printf("Step %d done\n", stepi)
    }
}

func (brewer *Brewer) Abort() {
    brewer.abort_mutex.Lock()
    defer brewer.abort_mutex.Unlock()

    select {
    case <- brewer.abort:
    default:
        close(brewer.abort)
    }
}

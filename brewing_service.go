package main

import (
    "strconv"
    "sync/atomic"
    "sync"
)

type Recipe []map[string]int

type Brewer struct {
    is_brewing int32
    task_id int
    recipe Recipe
}

type BrewingService interface {
  Brew(recipe Recipe)
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

    for _,step := range recipe {
        var stepWg sync.WaitGroup
        for idStr,time := range step {
          var id, _ = strconv.ParseInt(idStr, 10, 64)
          stepWg.Add(1)
          go ServePump(id, time, stepWg)
        }

        stepWg.Wait()
    }
}

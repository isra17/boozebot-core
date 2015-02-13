package main

import (
    "strconv"
    "sync/atomic"
    "sync"
    "fmt"
    "time"
    "github.com/isra17/boozebot-core/event"
)

type Recipe []map[string]int

type Brewer struct {
    is_brewing int32
    task_id int
    recipe Recipe

    abort event.Event
    pause event.Event
    resume event.Event
}

func MakeBrewer() Brewer {
  brewer := Brewer {
    abort: event.MakeEvent(),
    pause: event.MakeEvent(),
    resume: event.MakeEvent(),
  }

  return brewer
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
    brewer.abort.Reset()
    brewer.pause.Reset()
    brewer.resume.Reset()

    for stepi,step := range recipe {
        if !brewer.processStep(step) {
          fmt.Printf("Abort\n")
          break
        }
        fmt.Printf("Step %d done\n", stepi)
    }
}

func updateStep(step map[string]int, elapsed int) map[string]int {
  for id,time := range step {
    step[id] -= elapsed
    if(time <= elapsed) {
      delete(step, id)
    }
  }

  return step
}

func (brewer *Brewer) waitPause() bool {
  brewer.resume.Reset()
  select {
  case <-brewer.resume.Wait():
    return true
  case <-brewer.abort.Wait():
    return false
  case <-time.After(time.Second * 10):
    return false
  }
}

func (brewer *Brewer) processStep(step map[string]int) bool {
  for len(step) > 0 {
    startAt := time.Now().UnixNano()

    stepWg := sync.WaitGroup{}
    stepWg.Add(len(step))

    for idStr,time := range step {
      var id, _ = strconv.ParseInt(idStr, 10, 64)
      go ServePump(id, time, &stepWg, brewer)
    }

    stepWg.Wait()

    select {
      case <- brewer.pause.Wait():
        elapsed := int((time.Now().UnixNano() - startAt) / 1000000)
        fmt.Printf("Elapsed: %d\n", elapsed)
        step = updateStep(step, elapsed)

        if !brewer.waitPause() {
          return false
        }

        brewer.pause.Reset()
        fmt.Println("Resume")
      case <- brewer.abort.Wait():
        return false
      default:
        return true
    }
  }

  return true
}

func (brewer *Brewer) Abort() {
  brewer.abort.Signal()
}

func (brewer *Brewer) Pause() {
  brewer.pause.Signal()
}

func (brewer *Brewer) Resume() {
  brewer.resume.Signal()
}

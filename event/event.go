package event

import "sync"

type empty struct {}

type Event struct {
  m sync.Mutex
  c chan empty
}

func MakeEvent() Event {
  e := Event{m: sync.Mutex{}, c: make(chan empty)}
  return e
}

func (e *Event) Signal() {
  e.m.Lock()
  defer e.m.Unlock()

  e.unsafe_signal()
}

func (e *Event) unsafe_signal() {
  select {
  case <- e.c:
  default:
    close(e.c)
  }
}

func (e *Event) Reset() {
  e.m.Lock()
  defer e.m.Unlock()

  e.unsafe_signal()
  e.c = make(chan empty)
}

func (e *Event) Wait() chan empty {
  return e.c
}

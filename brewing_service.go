package main

import (
    "time"
    "sync/atomic"
)

type Recipe []map[string]int

type Brewer struct {
    is_brewing int32
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
    defer brewer.Unlock()

    time.Sleep(3000 * time.Millisecond)

    return true;
}

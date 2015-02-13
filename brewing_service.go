package drinkifycore;

import (
    "sync/atomic"
)

type Recipe []map[string]int

type BrewingState struct {
    is_brewing int32
}

type Brewer interface {
  Brew(recipe Recipe)
}

func (brewer *Brewer) Lock() {
    return atomic.CompareAndSwapInt32(&brewer.is_brewing, 0, 1)
}

func (brewer *Brewer) Unlock() {
    brewer.is_brewing = 0
}

func (brewer *Brewer) Brew() bool {
    if !Lock() { return false }
    defer Unlock()

    time.Sleep(3000 * time.Millisecond)

    return true;
}

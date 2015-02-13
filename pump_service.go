package main

import (
    "time"
    "sync"
    "fmt"
)

func ServePump(id int64, ms int, wg *sync.WaitGroup, abort <-chan struct{}) {
    defer wg.Done()
    select {
    case <-time.After(time.Duration(ms) * time.Millisecond):
        fmt.Printf("Pump %d done after %d ms\n", id, ms)
    case <-abort:
        fmt.Printf("Pump %d aborted", id)
    }
}

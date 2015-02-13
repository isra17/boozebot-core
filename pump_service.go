package main

import (
    "time"
    "sync"
    "fmt"
)

func ServePump(id int64, ms int, wg *sync.WaitGroup) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
    wg.Done()

    fmt.Printf("Pump %d done after %d ms\n", id, ms)
}

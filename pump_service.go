package main

import (
    "time"
    "sync"
)

func ServePump(id int64, ms float64, wg sync.WaitGroup) {
    time.Sleep(ms * time.Millisecond)
    wg.Done()
}

package main

import (
    "time"
    "sync"
    "fmt"
)

func ServePump(id int64, ms int, wg *sync.WaitGroup, brewer *Brewer) {
    defer wg.Done()
    physicalPump := physicalPumps[id - 1]
    physicalPump.DigitalWrite(1)
    select {
    case <-time.After(time.Duration(ms) * time.Millisecond):
        fmt.Printf("Pump %d done after %d ms\n", id, ms)
    case <-brewer.abort.Wait():
        fmt.Printf("Pump %d aborted\n", id)
    case <-brewer.pause.Wait():
        fmt.Printf("Pump %d paused\n", id)
    }
    physicalPump.DigitalWrite(0)
}

package main

import (
    "time"
    "encoding/json"
    "fmt"
    "net/http"
    "sync/atomic"
)

type Recipe map[string]int

type Mutex struct {
    isBoozing int32
}
var boozing = &Mutex{}

func brew(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world!\n")
    decoder := json.NewDecoder(r.Body)
    var boozeRequest []Recipe
    err := decoder.Decode(&boozeRequest)
    if(err != nil) {
        fmt.Fprintf(w, "Wrong request format\n" + err.Error())
        return
    }
    if(atomic.CompareAndSwapInt32(&boozing.isBoozing, 0, 1)) {
        time.Sleep(3000 * time.Millisecond)
        fmt.Fprintf(w, "%+v\n", boozeRequest)
        boozing.isBoozing = 0
    } else {
        fmt.Fprintf(w, "Already boozing\n")
    }
}

func main() {
    http.HandleFunc("/brew", brew)
    http.ListenAndServe(":6543", nil)
}

package main

import (
    "time"
    "encoding/json"
    "fmt"
    "net/http"
    "sync/atomic"
    "brewing_service"
)

var brewer = Brewer

func brew(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world!\n")
    decoder := json.NewDecoder(r.Body)
    var brewRequest Recipe
    err := decoder.Decode(&brewRequest)
    if(err != nil) {
        fmt.Fprintf(w, "Wrong request format\n" + err.Error())
        return
    }
    if brewer.Brew(brewRequest) {
        fmt.Fprintf(w, "%+v\n", brewRequest)
    } else {
        fmt.Fprintf(w, "Already boozing\n")
    }
}

func main() {
    http.HandleFunc("/brew", brew)
    http.ListenAndServe(":6543", nil)
}

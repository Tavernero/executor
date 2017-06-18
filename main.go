package main

import (
    "log"
    "fmt"
    "net/http"
)

func main() {

    log.Println("main start")

    TaskQueue = make(chan Task, MaxQueue)

    configuration := Configuration{MaxWorkers:MaxWorker,Function:"create"}
    dispatcher := NewDispatcher(configuration)
    dispatcher.Run()


    http.HandleFunc("/task/", postHandlerCreateTask)
    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        fmt.Println("starting listening for payload messages")
    } else {
        fmt.Errorf("an error occured while starting payload server %s", err.Error())
    }
}

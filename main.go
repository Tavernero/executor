package main

import (
    "log"
)

// XXX : how to defer close db connections
// XXX : getoption lib == https://golang.org/pkg/flag/

func main() {

    TaskQueue = make(chan Task, MaxQueue)

    go launchHttpServer()

    log.Println("main start")

    configuration := Configuration{
        MaxWorkers:MaxWorker,
        MaxDatabaseWorkers:MaxDatabaseWorker,
        Function:"create" }

    dispatcher := NewDispatcher(configuration)

    dispatcher.Run()
}

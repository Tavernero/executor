package main

import (
    "log"
)

// XXX : how to defer close db connections
// XXX : getoption lib == https://golang.org/pkg/flag/

func main() {

    log.Println("Main start")

    configuration := Configuration{
        MaxWorkers:MaxWorker,
        MaxQueue:MaxQueue,
        Function:"create" }

    dispatcher := NewDispatcher(configuration)

    dispatcher.Run()
}

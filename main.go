package main

import (
    "log"
)

func main() {

    log.Println("Main start")

    configuration := Configuration{
        MaxWorkers:MaxWorker,
        MaxQueue:MaxQueue,
        Function:"create" }

    dispatcher := NewDispatcher(configuration)

    dispatcher.Run()
}

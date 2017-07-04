package main

import (
    "log"
    "strconv"
    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"
)

// Dispatcher object
type Dispatcher struct {
    Definition      Definition
    Configuration   Configuration
    WorkerPool      chan chan int64
    IdQueue         chan int64

    connector       *gorp.DbMap
    quit            chan bool
}

// Create a new dispatcher
func NewDispatcher(configuration Configuration) *Dispatcher {
    pool := make(chan chan int64, configuration.MaxWorkers)
    queue := make(chan int64, configuration.MaxQueue)
    return &Dispatcher{
        Configuration: configuration,
        WorkerPool: pool,
        IdQueue: queue }
}

// Launch the dispatcher process
func (d *Dispatcher) Run() {

    // retrieve a gorp dbmap
    d.connector = initDb()
    // XXX : defer d.connector.Db.Close()

    // retrieve the steps for this function
    d.readSteps()

    // starting n number of workers
    for i := 0; i < d.Configuration.MaxWorkers; i++ {
        worker := NewWorker(d.Configuration.Function,d.WorkerPool,d.Definition)
        worker.Start()
    }

    // launch a first read on database data task
    go d.readTaskIds()

    // launch the listener for database events
    go d.initializeListenerAndListen()

    // launch the dispatch
    d.dispatch()
}

// Dispatch each task to a free worker
func (d *Dispatcher) dispatch() {

    log.Println("Worker queue dispatcher started...")

    for {
        select {
            case taskId := <-d.IdQueue:

                log.Printf("Dispatch to taskChannel with ID : " + strconv.Itoa( int(taskId) ) )

                // try to obtain a worker task channel that is available.
                // this will block until a worker is idle
                taskChannel := <-d.WorkerPool

                // dispatch the task to the worker task channel
                taskChannel <- taskId

            case <-d.quit:
                // we have received a signal to stop

                // XXX : how to stop workers correctly
        }
    }
}

// XXX : how to improve this part ?
// Stop signals the worker to stop listening for work requests.
func (d *Dispatcher) Stop() {
    go func() {
        d.quit <- true
    }()
}

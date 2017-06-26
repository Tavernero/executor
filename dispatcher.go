package main

import (
    "fmt"
    "log"
    "strconv"
    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"
)

// Dispatcher object
type Dispatcher struct {
    Configuration   Configuration
    WorkerPool      chan chan Task
    TaskQueue       chan Task

    quit            chan bool
    Steps           []DispatcherStep
    connector       *gorp.DbMap
}

// Create a new dispatcher
func NewDispatcher(conf Configuration) *Dispatcher {
    pool := make(chan chan Task, conf.MaxWorkers)
    queue := make(chan Task, conf.MaxQueue)
    return &Dispatcher{
        WorkerPool: pool,
        TaskQueue: queue,
        Configuration: conf }
}

// Launch the dispatcher process
func (d *Dispatcher) Run() {

    // retrieve a gorp dbmap
    d.connector = initDb()
    // XXX : defer d.connector.Db.Close()

    // retrieving steps from db
    _, err := d.connector.Select(
        &d.Steps,
        "select * from task_step where function = :function order by index asc",
        map[string]interface{}{"function":d.Configuration.Function} )

    if err != nil {
        log.Fatalln("Select failed", err)
    }

    log.Println("All task steps :")

    for x, p := range d.Steps {
        log.Printf("  %d  :  %v", x, p)
    }

    log.Println("===============")

    // starting n number of workers
    for i := 0; i < d.Configuration.MaxWorkers; i++ {
        worker := NewWorker(d.WorkerPool,d.Steps)
        worker.Start()
    }

    // launch a first read on database data task
    go d.firstRead()

    // launch the listener for database events
    go d.initializeListenerAndLaunch()

    // launch the dispatch
    d.dispatch()
}

// Dispatch each task to a free worker
func (d *Dispatcher) dispatch() {
    fmt.Println("Worker queue dispatcher started...")
    for {
        select {
            case task := <-d.TaskQueue:

                log.Printf("Dispatch to taskChannel with ID : " + strconv.Itoa( task.ID ) )

//                // a task request has been receive
//                go func(task Task) {

                // try to obtain a worker task channel that is available.
                // this will block until a worker is idle
                taskChannel := <-d.WorkerPool

                // dispatch the task to the worker task channel
                taskChannel <- task

//                }(task)

            case <-d.quit:
                // we have received a signal to stop

                // XXX : how to stop workers correctly

                return
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

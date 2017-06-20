package main

import (
    "fmt"
    "log"

    "database/sql"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"
)

// ================================================= //
// ================================================= //

// Worker represents the worker that executes the task
type Worker struct {
    WorkerPool  chan chan Task
    TaskChannel chan Task
    quit        chan bool

    Steps []DispatcherStep

//    connector gorp.DbMap
    connector *gorp.DbMap
//    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
//    return dbmap
}

// Create a new worker
func NewWorker(workerPool chan chan Task,steps []DispatcherStep) Worker {
    return Worker{
        WorkerPool: workerPool,
        TaskChannel: make(chan Task),
        quit: make(chan bool),
        Steps: steps}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {

    // BEGIN
    // Retrieve a new connection for each worker instance
    db, err := sql.Open("postgres", ConnectionConfiguration)
    if err != nil {
        log.Fatalln("sql.Open failed ...", err )
        fmt.Println("toto is back no???")
        panic(err)
    }
    w.connector = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
    // END

    go func() {
        for {
            // register the current worker into the worker queue.
            w.WorkerPool <- w.TaskChannel

            // read from the channel
            select {
                case task := <-w.TaskChannel:
                    // we have received a work request.
                    if err := task.WorkingOn(); err != nil {
                        log.Printf("Error while working on task: %s", err.Error())
                    }

                case <-w.quit:
                    // we have received a signal to stop
                    return
            }
        }
    }()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}

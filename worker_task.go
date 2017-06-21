package main

import (
    "log"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"

    "fmt"
    "math/rand"
    "strconv"
    "time"
)

// ================================================= //
// ================================================= //

// Worker represents the worker that executes the task
type Worker struct {
    WorkerPool  chan chan Task
    TaskChannel chan Task
    quit        chan bool
    Steps       []DispatcherStep
    connector   *gorp.DbMap
}

// Create a new worker
func NewWorker(workerPool chan chan Task,steps []DispatcherStep) Worker {
    return Worker{
        WorkerPool: workerPool,
        TaskChannel: make(chan Task),
        quit: make(chan bool),
        Steps: steps }
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {

    // retrieve a gorp dbmap
    w.connector = initDb()
    // XXX : defer d.connector.Db.Close()

    go func() {
        for {
            // register the current worker into the worker queue.
            w.WorkerPool <- w.TaskChannel

            // read from the channel
            select {
                case task := <-w.TaskChannel:
                    log.Println("entry to task channels")

                    // we have received a work request.
                    if err := w.Action(task); err != nil {
                        log.Printf("Error while working on task: %s", err.Error())
                    }
//                    // we have received a work request.
//                    if err := task.WorkingOn(); err != nil {
//                        log.Printf("Error while working on task: %s", err.Error())
//                    }

                case <-w.quit:
                    // we have received a signal to stop
                    return
            }
        }
    }()
}


func (w Worker) Action(task Task) error {


    randomer := rand.Intn(2000)
    time.Sleep(time.Duration(randomer) * time.Millisecond)
    fmt.Println("work done in " + strconv.Itoa( randomer ) + " ms for '" + strconv.Itoa( task.ID ) + "'" )

    task.Status = "doing"

    res, err := w.connector.Update(&task)

    if err != nil {
        log.Fatalln("update fialed", err)
    }

    log.Println("Rows updated:", res)





    return nil
}


// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}

















































































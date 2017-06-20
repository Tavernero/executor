package main

import (
    "fmt"
    "log"

    "encoding/json"
    "time"
    "strconv"

    "database/sql"
    "github.com/lib/pq"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"

    "bytes"
)

// Dispatcher object
type Dispatcher struct {
    Configuration Configuration
    WorkerPool chan chan Task
    DatabaseWorkerPool chan chan *gorp.DbMap

    connector *gorp.DbMap

    Steps []DispatcherStep
}

// Create a new dispatcher
func NewDispatcher(conf Configuration) *Dispatcher {

//    // get the database informations
//    var dbmapper = initDb()
//
//    // retrieve all task
//    var tasks []Task
//
////    _, err := dbmapper.Select(&tasks, "select id, function, name, step, status, retry from task order by id asc")
//    _, err := dbmapper.Select(&tasks, "select * from task order by id asc")
//
//    if err != nil {
//        log.Fatalln("Select failed", err)
//    }
//
//    log.Println("All rows:")
//
//    for x, p := range tasks {
//        log.Printf("  %d : %v\n", x, p)
//    }

    pool := make(chan chan Task, conf.MaxWorkers)
    databasePool := make(chan chan *gorp.DbMap, conf.MaxWorkers)
    return &Dispatcher{WorkerPool: pool, DatabaseWorkerPool: databasePool, Configuration: conf}
}

// Launch the dispatcher process
func (d *Dispatcher) Run() {

    // retrieve a gorp dbmap
    d.connector = initDb()
    defer d.connector.Db.Close()

    // retrieve from "function" name a configuration for the task steps/calls
//    d.Steps = []DispatcherStep{
//        DispatcherStep{Name:"starting",Url:"https://api.com/starting"},
//        DispatcherStep{Name:"onServer",Url:"https://api.com/onServer"},
//        DispatcherStep{Name:"onInterne",Url:"https://api.com/onInterne"},
//        DispatcherStep{Name:"ending",Url:"https://api.com/ending"}}

    // get the database informations
    var dbmapper = initDb()

    // retrieve all task
    _, err := dbmapper.Select(&d.Steps, "select * from task_step order by id asc")

    if err != nil {
        log.Fatalln("Select failed", err)
    }

    log.Println("All task steps :")

    for x, p := range d.Steps {
        log.Printf("  %d  :  %v  \n", x, p)
    }

//    // update a row
//    p2.Title = "Go 1.2 is better than ever"
//    count, err = dbmap.Update(&p2)
//    checkErr(err, "Update failed")
//    log.Println("Rows updated:", count)
//
//    // confirm count is zero
//    count, err = dbmap.SelectInt("select count(*) from posts")
//    checkErr(err, "select count(*) failed")
//    log.Println("Row count - should be zero:", count)

    // starting n number of workers
    for i := 0; i < d.Configuration.MaxWorkers; i++ {
        worker := NewWorker(d.WorkerPool,d.Steps)
        worker.Start()
    }

    // launch the dispatch
    go d.dispatch()

    // launch the listener for database events
    go d.initializeListenerAndLaunch()
}

// Dispatch each task to a free worker
func (d *Dispatcher) dispatch() {
    fmt.Println("Worker queue dispatcher started...")
    for {
        select {
            case task := <-TaskQueue:

                log.Printf("Dispatch to taskChannel with ID : " + strconv.Itoa( task.ID ) )

                // a task request has been receive
                go func(task Task) {
                    // try to obtain a worker task channel that is available.
                    // this will block until a worker is idle
                    taskChannel := <-d.WorkerPool

                    // dispatch the task to the worker task channel
                    taskChannel <- task
                }(task)
        }
    }
}




// ================================================= //
// ================================================= //


// prepare the listener data and launch it
func (d *Dispatcher) initializeListenerAndLaunch() {

    _, err := sql.Open("postgres", ConnectionConfiguration)

    if err != nil {
        panic(err)
    }

    reportProblem := func(ev pq.ListenerEventType, err error) {
        if err != nil {
            fmt.Println(err.Error())
        }
    }

    listener := pq.NewListener(ConnectionConfiguration, 10*time.Second, time.Minute, reportProblem)

    err = listener.Listen("events")

    if err != nil {
        panic(err)
    }

    fmt.Println("Start monitoring PostgreSQL...")

    for {
        d.waitForNotificationFromListener(listener)
    }
}

// listening to the event bus of the database and do some actions
func (d *Dispatcher) waitForNotificationFromListener(l *pq.Listener) {
    for {
        select {
            case n := <-l.Notify:
                fmt.Println("Received data from channel [", n.Channel, "] :")
                // Prepare notification payload for pretty print
                var prettyJSON bytes.Buffer
                err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
                if err != nil {
                    fmt.Println("Error processing JSON: ", err)
                    return
                }
                fmt.Println(string(prettyJSON.Bytes()))



                // get the database informations
                var dbmapper = initDb()

                // retrieve all task
                var tasks []Task

                _, err = dbmapper.Select(&tasks, "select * from task order by id asc")

                if err != nil {
                    log.Fatalln("Select failed", err)
                }

                log.Println("All rows:")

                for x, task := range tasks {
                    TaskQueue <- task

                    log.Printf("  %d : %v\n", x, task)
                }

                return

            case <-time.After(90 * time.Second):
                fmt.Println("Received no events for 90 seconds, checking connection")

                go func() {
                    l.Ping()
                }()

                // get the database informations
                var dbmapper = initDb()

                // retrieve all task
                var tasks []Task

                _, err := dbmapper.Select(&tasks, "select * from task order by id asc")

                if err != nil {
                    log.Fatalln("Select failed", err)
                }

                log.Println("All rows:")

                for x, task := range tasks {
                    TaskQueue <- task

                    log.Printf("  %d : %v\n", x, task)
                }

                return
        }
    }
}

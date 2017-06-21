package main

import (
    "fmt"
    "log"

    "encoding/json"
    "time"

    "database/sql"
    "github.com/lib/pq"
)

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


type DeltaTask struct {
    Table string
    Action string
    Data Task
}

// listening to the event bus of the database and do some actions
func (d *Dispatcher) waitForNotificationFromListener(l *pq.Listener) {
    for {
        select {
            case n := <-l.Notify:
                fmt.Println("Received data from channel [", n.Channel, "] :")

                var task DeltaTask

                err := json.Unmarshal([]byte(n.Extra), &task)

                if err != nil {
                    fmt.Println("error:",err)
                }

                fmt.Printf("%+v", task)

                TaskQueue <- task.Data

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

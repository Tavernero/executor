package main

import (
    "log"
    "time"
    "encoding/json"
    "database/sql"
    "github.com/lib/pq"
)

func (d *Dispatcher) readSteps() {

    log.Println("Get the dispatcher definition")

    err := d.connector.SelectOne(
        &d.Definition,
        "select * from definition where function = :function",
        map[string]interface{}{"function":d.Configuration.Function} )

    if err != nil {
        log.Fatalln("Select failed", err)
    }
}

func (d *Dispatcher) readTaskIds() {

    log.Println("Read all task IDs")

    var taskIds []int64

    _, err := d.connector.Select(
        &taskIds,
        "select id from task where status = :status and function = :function and retry > 0 and todo_date <= now() order by id asc",
        map[string]interface{}{"status":"todo","function":d.Configuration.Function} )

    if err != nil {
        log.Fatalln("Select failed", err)
    }

    log.Println("All rows:")

    for x, id := range taskIds {

        d.IdQueue <- id

        log.Printf("  %d : %v\n", x, id)
    }
}

func (d *Dispatcher) initializeListenerAndListen() {

    _, err := sql.Open("postgres", ConnectionConfiguration)

    if err != nil {
        panic(err)
    }

    reportProblem := func(ev pq.ListenerEventType, err error) {
        if err != nil {
            log.Println(err.Error())
        }
    }

    listener := pq.NewListener(ConnectionConfiguration, 10*time.Second, time.Minute, reportProblem)

    err = listener.Listen("events_task_" + d.Configuration.Function)

    if err != nil {
        panic(err)
    }

    log.Println("Start monitoring PostgreSQL...")

    for {
        d.waitForNotification(listener)
    }
}

type DatabaseNotification struct {
    Table       string
    Action      string
    Function    string
    ID          int64
}

func (d *Dispatcher) waitForNotification(l *pq.Listener) {
    for {
        select {
            case n := <-l.Notify:

                var notification DatabaseNotification

                err := json.Unmarshal([]byte(n.Extra), &notification)

                if err != nil {
                    log.Println("error:",err)
                }

                log.Println("Received data from channel [", n.Channel, "] :")

                log.Printf("%+v \n", notification)

                d.IdQueue <- notification.ID

                log.Println("Data send in task queue")

            case <-time.After(90 * time.Second):

                log.Println("Received no events for 90 seconds, checking connection")

                go l.Ping()

                log.Println("Retreieve ids")

                go d.readTaskIds()
        }
    }
}

func (w *Worker) readOneTask(id int64) (task Task, err error) {

    log.Println("Read one task")

    err = w.connector.SelectOne(
        &task,
        "select * from task where status = :status and function = :function and id = :id and retry > 0 and todo_date <= now() limit 1",
        map[string]interface{}{"status":"todo","function":w.Function,"id":id} )

    if err != nil {
        log.Fatalln("Select failed", err)
    }

    return task, err
}


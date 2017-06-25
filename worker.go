package main

import (
    "log"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"

    "fmt"
//    "math/rand"
    "strconv"
//    "time"

    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
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

                    log.Printf("Entry to taskChannel with ID : " + strconv.Itoa( task.ID ) + "\n")

                    // we have received a work request.
                    if err := w.Action(task); err != nil {
                        log.Printf("Error while working on task: %s", err.Error())
                    }

                    log.Println(".")
                    log.Println(".")
                    log.Println("ENDJOB")
                    log.Println(".")
                    log.Println(".")
//                    // we have received a work request.
//                    if err :. task.WorkingOn(); err !. nil {
//                        log.Printf("Error while working on task: %s", err.Error())
//                    }

                case <-w.quit:
                    // we have received a signal to stop
                    return
            }
        }
    }()
}

type HttpOut struct {
    Name string
    Arguments JsonB
    Buffer JsonB
}




func (w Worker) Action(task Task) error {

    // vars
    var ending = false

    // update task status
    task.Status = "doing"
    // task.LastUpdate = time.Now()

    // exception for ending steps
    if task.Step == "ending" {
        task.Status = "done"
        // task.DoneDate = time.Now()
        ending = true
    }

    res, err := w.connector.Update(&task)

    if err != nil {
        log.Fatalln("Error while updating the task for status lock", err)
    }

    if res != 1 {
        log.Fatalln("Error while updating the task status")
    }

    // exception for ending step
    if ending {
        return nil
    }

    // store step informations
    var stepId = -1
    var stepData DispatcherStep

    // retrieve step data
    for i, s := range w.Steps {
        log.Println(" id . " + strconv.Itoa(i) + " step . " + s.Name )

        if s.Name == task.Step {
            stepId = i
            stepData = s
        }
    }

    // no associated step found, error
    if stepId == -1 {
        log.Fatalln("Error while finding the good task step informations")
    }

    log.Println("Working on task " + strconv.Itoa( task.ID ) + "/" + task.Function + " on step: " + task.Step )

    // initialize the http client
    httpclient := http.Client{}



    var dataOut = HttpOut{
        Name: task.Name,
        Arguments: task.Arguments,
        Buffer: task.Buffer}

    jsonValue, _ := json.Marshal(dataOut)

    req, err := http.NewRequest("POST", stepData.Url, bytes.NewBuffer(jsonValue))

    req.Header.Set("X-Custom-Header", "myvalue")

    req.Header.Set("Content-Type", "application/json")

    log.Println("===============================")
    log.Printf("Post data request was '", string(jsonValue), "'")

    resp, err := httpclient.Do(req)

    if err != nil {
        fmt.Println("error while request", err)
        panic(err)
    }

    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)

    fmt.Println("response Headers:", resp.Header)

    body, _ := ioutil.ReadAll(resp.Body)

    fmt.Println("response Body:", string(body))

    log.Println("-------------------------------")

    var statusCode = resp.StatusCode

    if statusCode != 200 {

        // need to match errors states


        task.Status = "error"


        res, err = w.connector.Update(&task)

        if err != nil {
            log.Fatalln("update fialed", err)
        }

        log.Println("Rows updated:", res)



        return nil
    }

    // 200 ok, 

    task.Status = "todo"
    // task.LastUpdate = time.Now()


    var dataJson JsonB

    err = json.Unmarshal(body, &dataJson)
    //err := json.Unmarshal([]byte(n.Extra), &notification)

    if err != nil {
        fmt.Println("error:",err)
    }

    //task.Buffer[ task.Step ] = string(body) // maybe decaps to recaps into go struct and not only string
    task.Buffer[ task.Step ] = dataJson

    res, err = w.connector.Update(&task)

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














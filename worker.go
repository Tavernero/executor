package main

import (
    "log"
    "strconv"
    "net/http"
    "io/ioutil"
    "bytes"
    "time"
    "encoding/json"
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
                    log.Println("ENDJOB")
                    log.Println(".")

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



























//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//  200  = next                                                         //
//  302  = next to '...' step or/and next in '...' interval of seconds  //
//  420  = cancelled                                                    //
//  520  = problem                                                      //
// other = error ( 5XX : auto-retry )                                   //
//______________________________________________________________________//

// Response http from each API response
type HttpResponse struct {
    Buffer      *JsonB
    Interval    *int
    Step        *string
    Comment     *string
}

func (w Worker) Action(task Task) error {

    // Little sleeper for best test works
    time.Sleep( time.Second )

    // vars
    var ending = false

    // update task status
    task.Status = "doing"
    //task.LastUpdate = time.Now()

    log.Println(" --=-=-=-=-=-=-=-=-=-=-=------------=-=-=-=-=-=-=-=-=-=-=-=-=--")
    log.Println( task )
    log.Println(" --=-=-=-=-=-=-=-=-=-=-=------------=-=-=-=-=-=-=-=-=-=-=-=-=--")

    // exception for ending steps
    if task.Step == "ending" {
        task.Status = "done"
        // task.DoneDate = time.Now()
        ending = true
    }

    res, err := w.connector.Update(&task)

    if err != nil {
        log.Fatalln("Error while updating the task for status lock : ", err)
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

//        log.Println("Steps => id . " + strconv.Itoa(i) + " step . " + s.Name )

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









    // Do the http call to retrieve API data/informations
    httpResponse, statusCode, err := w.CallHttp(task, stepData)


//    log.Println("httpResponse : ")
//    log.Println( httpResponse )
//    log.Println("statusCode : ")
//    log.Println( statusCode )
//    log.Println("err : ")
//    log.Println( err )

    if err != nil {

        // XXX : need to re-read this part, maybe badly write
        log.Fatalln("Error while do the http call", err)

        var comment = "Error while doing the http call"

        // No http response, we "emulate" a 500 error with a fake http response
        httpResponse = HttpResponse{
            Buffer: &task.Buffer,
            Comment: &comment }

        statusCode = 500
    }

    // Switch on each status code
    switch statusCode {

        // 200 = next
        case 200:

            // Which step next ?
            var nextStep = w.Steps[ stepId + 1 ]

            if nextStep.ID == 0 {
                log.Fatalln("Error while retrieve the next step")
            }

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
//  200  = next                                                         //
//  302  = next to '...' step or/and next in '...' interval of seconds  //
//  420  = cancelled                                                    //
//  520  = problem                                                      //
// other = error ( 5XX : auto-retry )                                   //
//______________________________________________________________________//

            // Update the task informations
            task.Status = "todo"
            task.Step = nextStep.Name
            task.Buffer = *httpResponse.Buffer
//            task.LastUpdate = time.Now()





        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
        //  301  = next to '...' step or/and next in '...' interval of seconds  //
        //______________________________________________________________________//
        case 301:

            var founded = false

            // An interval is setup to TTL the next execution
            if httpResponse.Interval != nil {

                var newInterval = *httpResponse.Interval // in seconds

                if newInterval <= 0 {
                    log.Println("Interval not founded")
                } else {
                    var todoDateData = *task.TodoDate
                    log.Println("TodoDate : " + todoDateData.String() )
                    log.Println("Interval founded : " + strconv.Itoa( newInterval ) )
                    var newTodoDate = todoDateData.Add( time.Duration( newInterval ) * time.Second )
                    log.Println("New interval computed : " + newTodoDate.String() )
                    task.TodoDate = &newTodoDate
                    log.Println("The todoDate updated to : '" + newTodoDate.String() + "' with interval '" + strconv.Itoa( newInterval ) + "'")
                    founded = true
                }
            }

            // A new next step definition
            if httpResponse.Step != nil {

                var askedStep = *httpResponse.Step

                // flag founded
                var stepFounded = false

                // loop on each loaded steps data
                for _, s := range w.Steps {
                    // found the asked overwritted step
                    if s.Name == askedStep {
                        stepFounded = true
                    }
                }

                if stepFounded {
                    log.Println("Asked step founded : '" + askedStep + "'")
                    founded = true
                } else {
                    log.Println("Asked step error : '" + askedStep + "'")
                }
            }

            // Update the status data
            if founded {
                task.Status = "todo"
            } else {
                task.Status = "error"
            }





        //‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾//
        //  420  = cancelled                                                    //
        //______________________________________________________________________//
        case 420:

            // A comment was setup ???
            if httpResponse.Comment != nil {

                var comment = *httpResponse.Comment // string

                if comment != "" {
                    task.Comment = comment
                }
            }

            // Setup the status
            task.Status = "cancelled"




        // 520 = problem
        case 520:


            // 200 ok,

            // retrieve the next step informations
            nextStep := w.Steps[ stepId + 1 ]

            if nextStep.ID == 0 {
                log.Fatalln("Imposisble to found the next step")
            }

            // Decoding the return buffer informations
            var dataJson JsonB

//            err = json.Unmarshal(body, &dataJson)
//            //err := json.Unmarshal([]byte(n.Extra), &notification)
//
//            if err != nil {
//                fmt.Println("error:",err)
//            }

            // Update the task informations
            task.Status = "todo"
            task.Step = nextStep.Name
            task.Buffer = dataJson
            // task.LastUpdate = time.Now()

            // Do request on the database
            res, err = w.connector.Update(&task)

            if err != nil {
                log.Fatalln("update fialed", err)
            }

            log.Println("Rows updated:", res)

            return nil

        // other = error ( 5XX : auto-retry )
        default:

            // Default status
            status := "error"

            // Exception for 5XX status code, auto retry
            if ( statusCode / 100 ) == 5 {
                status = "todo"
                // retry -= retry
            }

            // Update the task informations
            task.Status = status
            task.Buffer = *httpResponse.Buffer
//            task.LastUpdate = time.Now()

    }

    // Update the last update date
    var timeNow = time.Now()
    task.LastUpdate = &timeNow

    // Update the retry counter
    task.Retry = task.Retry - 1

    // Do request on the database
    num, err := w.connector.Update(&task)

    if err != nil {
        log.Fatalln("Error while update on the database the task", err)
    }

    if num > 1 {
        log.Fatalln("Error while updating the task, more than one row modified")
    }

    if num < 1 {
        log.Fatalln("Error while updating the task, no row modified")
    }

    log.Println("Task updated")

    return nil
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}
















































func (w Worker) CallHttp(task Task, step DispatcherStep) ( httpResponse HttpResponse, statusCode int, err error ) {

    // Initialize the http client
    httpclient := http.Client{}

    // Http call data
    var dataOut = HttpOut{
        Name: task.Name,
        Arguments: task.Arguments,
        Buffer: task.Buffer}

    // Encode the http call data
    jsonValue, err := json.Marshal(dataOut)

    if err != nil {
        log.Fatalln("Error while encode the http call data", err)

        return httpResponse, statusCode, err
    }

    // Create the http request
    req, err := http.NewRequest("POST", step.Url, bytes.NewBuffer(jsonValue))

    if err != nil {
        log.Fatalln("Error while create the http resquest", err)

        return httpResponse, statusCode, err
    }

    // Set some headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Custom-Header", "my-custom-header")

    // Do the http call
    resp, err := httpclient.Do(req)

    if err != nil {
        log.Fatalln("Error while do the http call", err)

        return httpResponse, statusCode, err
    }

    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        log.Fatalln("Error while read the http body data", err)

        return httpResponse, statusCode, err
    }

    log.Println("===============================")
    log.Printf("Post data request was '%s'\n", string(jsonValue) )
    log.Println("Response Status:", resp.Status)
    log.Println("Response StatusCode:", resp.StatusCode)
    log.Println("Response Headers:", resp.Header)
    log.Println("Response Body:", string(body))
    log.Println("-------------------------------")

    // Decoding the returned body data
    err = json.Unmarshal(body, &httpResponse)

    if err != nil {
        log.Fatalln("Error while decoding the http response body", err)

        return httpResponse, statusCode, err
    }

    // Retrieve the statusCode data
    statusCode = resp.StatusCode

    log.Println("Http call work fine")

    return httpResponse, statusCode, nil
}

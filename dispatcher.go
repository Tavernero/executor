package main

import (
    "fmt"
    "log"

    "encoding/json"
    "io"
    "math/rand"
    "net/http"
    "time"
    "strconv"

    "database/sql"
    "github.com/lib/pq"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"

    "bytes"
)

// ================================================= //
// ================================================= //

var ConnectionConfiguration = "postgres://executor:totoTOTO89@641a3187-5896-49c9-af7d-d8bed8187f79.pdb.ovh.net:21684/executor"

//func getConnectionString() string {
//    host := getParamString("db.host", "")
//    port := getParamString("db.port", "3306")
//    user := getParamString("db.user", "")
//    pass := getParamString("db.password", "")
//    dbname := getParamString("db.name", "auction")
//    protocol := getParamString("db.protocol", "tcp")
//    dbargs := getParamString("dbargs", " ")
//
//    if strings.Trim(dbargs, " ") != "" {
//        dbargs = "?" + dbargs
//    } else {
//        dbargs = ""
//    }
//    return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
//        user, pass, protocol, host, port, dbname, dbargs)
//}

var dbmap = initDb()

func initDb() *gorp.DbMap {
    db, err := sql.Open("postgres", ConnectionConfiguration)
    if err != nil {
        log.Fatalln("sql.Open failed ...", err )
        panic(err)
    }
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}


//        // construct a gorp DbMap setting dialect to sqlite3
//        dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
//        defer dbmap.Db.Close()
//
//        // add a table, setting the table name to 'posts' and
//        // specifying that the Id property is an auto incrementing PK
//        dbmap.AddTableWithName(Car{}, "car").SetKeys(true, "ID")
//
//        // create the table. in a production system you'd generally
//        // use a migration tool, or create the tables via scripts
//        dbmap.CreateTablesIfNotExists()
//
//        var id = uuid.New()
//
//        dbmap.Insert(&Car{
//            ID: id,
//            Description: "Old Beater",
//            Color: "Brown",
//        })
//        var car *Car
//        dbmap.Get(car, id)



    return dbmap
}

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

// ================================================= //
// ================================================= //

// A buffered channel that we can send work requests on.
var TaskQueue chan Task

var (
    MaxWorker       = 20  //os.Getenv("MAX_WORKERS")
    MaxQueue        = 5 //os.Getenv("MAX_QUEUE")
    MaxLength int64 = 20480
)

// ================================================= //
// ================================================= //

// function called when we launch a task object
func (p *Task) WorkingOn() error {
    randomer := rand.Intn(2000)
    time.Sleep(time.Duration(randomer) * time.Millisecond)
    fmt.Println("work done in " + strconv.Itoa( randomer ) + " ms for '" + strconv.Itoa( p.ID ) + "'" )
    return nil
}

// task array from http request
type TaskCollection struct {
    Name string `json:"name"`
    Tasks []Task `json:"data"`
}

// handle the http post request on post localhost:8080/task/
func postHandlerCreateTask(w http.ResponseWriter, r *http.Request) {

    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // Read the body into a string for json decoding
    var content = &TaskCollection{}
    err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
    if err != nil {
        fmt.Errorf("an error occured while deserializing message")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    fmt.Println("Request received")

    // Go through each payload and queue items individually to be posted to S3
    for _, task := range content.Tasks {

        // Push the work onto the queue.
        TaskQueue <- task

        fmt.Println("Payload sent to workqueue : " + strconv.Itoa( task.ID ) )
    }

    w.WriteHeader(http.StatusOK)
}

// ================================================= //
// ================================================= //

// Task to work
type Task struct {
    ID int `db:"id, primarykey" json:"id"`
    Function string `db:"function" json:"function"`
    Name string `db:"name" json:"name"`
    Step string `db:"step" json:"step"`
    Status string `db:"status" json:"status"`
    Retry int `db:"retry" json:"retry"`
    Arguments string `db:"arguments" json:"arguments"`
    Buffer string `db:"buffer" json:"buffer"`
}

// Configuration for a specific task to work
type ConfigurationTask struct {
    ID int `db:"id, primarykey" json:"id"`
    Function string `db:"function" json:"function"`
    Status string `db:"status" json:"status"`
    Properties PropertyMap `db:"properties" json:"properties"`
}

// PropertyMap for catch JSONB from databases
type PropertyMap map[string]interface{}

//func (p PropertyMap) Value() (driver.Value, error) {
//    j, err := json.Marshal(p)
//    return j, err
//}
//
//func (p *PropertyMap) Scan(src interface{}) error {
//    source, ok := src.([]byte)
//    if !ok {
//        return errors.New("Type assertion .([]byte) failed.")
//    }
//
//    var i interface{}
//    err := json.Unmarshal(source, &i)
//    if err != nil {
//        return err
//    }
//
//    *p, ok = i.(map[string]interface{})
//    if !ok {
//        return errors.New("Type assertion .(map[string]interface{}) failed.")
//    }
//
//    return nil
//}

//create table "task_configuration" (
//    "id" bigserial primary key,
//    "function" character varying(255) not null,
//    "status" character varying(255) not null,
//    "properties" jsonb not null default '{}'
//);
//
//insert into "task_configuration" ( "function", "status", "properties" ) values
//( 'web/create', 'available', '{"sequence":[
//        {"step":"starting","url":"https://api.com/starting"},
//        {"step":"ending","url":"https://api.com/ending"}]}' );


// Dispatcher configuration object
type Configuration struct {
    Function string    // Function name where work
    MaxWorkers int    // A pool of workers channels that ardde registered with the dispatcher
}

// Dispatcher object
type Dispatcher struct {
    Configuration Configuration
    WorkerPool chan chan Task

    connector *gorp.DbMap

    Steps []DispatcherStep
}

// Step dispatcher object
type DispatcherStep struct {
    ID int
    Function string
    Name string
    Url string
}

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
    return &Dispatcher{WorkerPool: pool, Configuration: conf}
}

// Create a new worker
func NewWorker(workerPool chan chan Task,steps []DispatcherStep) Worker {
    return Worker{
        WorkerPool: workerPool,
        TaskChannel: make(chan Task),
        quit: make(chan bool),
        Steps: steps}
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

    log.Println("All rows:")

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

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {

    // BEGIN
    // Retrieve a new connection for each worker instance
    db, err := sql.Open("postgres", ConnectionConfiguration)
    if err != nil {
        log.Fatalln("sql.Open failed ...", err )
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

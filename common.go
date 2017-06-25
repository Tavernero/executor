package main

import (
    "log"

    "database/sql"

    "gopkg.in/gorp.v2"
    _ "github.com/lib/pq"
)

// ================================================= //
// ================================================= //

//var DEBUG_FLAG = false
var DEBUG_FLAG = true

// ================================================= //
// ================================================= //


// A buffered channel that we can send work requests on.
var TaskQueue chan Task

var (
    MaxWorker       = 20  //os.Getenv("MAX_WORKERS")
    MaxDatabaseWorker       = 20  //os.Getenv("MAX_WORKERS")
    MaxQueue        = 5 //os.Getenv("MAX_QUEUE")
    MaxLength int64 = 20480
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



    // add a table, setting the table name to 'task' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Task{}, "task").SetKeys(true, "ID")




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

// ================================================= //
// ================================================= //



// Dispatcher configuration object
type Configuration struct {
    Function string    // Function name where work
    MaxWorkers int    // A pool of workers channels that ardde registered with the dispatcher
    MaxDatabaseWorkers int    // A pool of workers channels that ardde registered with the dispatcher
}

// Step dispatcher object
type DispatcherStep struct {
    ID int
    Function string
    Index int
    Name string
    Url string
}

// Configuration for a specific task to work
//type ConfigurationTask struct {
//    ID int `db:"id, primarykey" json:"id"`
//    Function string `db:"function" json:"function"`
//    Status string `db:"status" json:"status"`
//    Properties PropertyMap `db:"properties" json:"properties"`
//}
//
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


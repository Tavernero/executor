package machine

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
    "time"

	"github.com/go-pg/pg"
	"github.com/m49n3t0/executor/models"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// to define the dispatcher parameters
type DispatcherParams struct {
	// robot configuration
	Function string
	// dispatcher configuration
	MaxWorker int64
	MaxQueue  int64
}

// dispatcher object
type Dispatcher struct {
	// function on what the robot work
	Function string
	// store datas from database
	Definitions map[int64]*models.Definition
	Endpoints   map[string]*models.Endpoint
	// manage the distribution of workflow
	WorkerPool chan chan string
	Queue      chan string
	MaxWorker  int64
	MaxQueue   int64
	// manage the quit process
	Signals chan os.Signal
	Quit    chan bool
	// database handler
	DB *pg.DB
	// logger handler
	Logger *logrus.Entry
}

// dispatcher creation handler
func NewDispatcher(params *DispatcherParams) (*Dispatcher, error) {


	// create a default logger
	var logger = logrus.New()

	// define the logger default level
	logger.Level = logrus.DebugLevel

	// define the logger default output
	logger.Out = os.Stdout

    logentry := logger.WithField("function", params.Function)

	// database configuration
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	// build database host address
	var address = host
	if port != "" {
		address = host + ":" + port
	}

	// pg database connector
	db := pg.Connect(&pg.Options{
		Addr:       address,
		User:       user,
		Password:   password,
		Database:   database,
		MaxRetries: 2,
	})

	// check connection
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		logentry.Println("Problem while check database connection")
		return nil, err
	}

	// build the query logger
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		// XXX : maybe only use UnformattedQuery option ( a debug flag ? )
		query, err := event.FormattedQuery()
		if err != nil {
			logentry.Panic(err)
		}

		logentry.WithField("sql_duration_ms", time.Since(event.StartTime)).Info(query)
	})


	// create the object
	dispatcher := &Dispatcher{
		Function:    params.Function,
		Definitions: make(map[int64]*models.Definition),
		Endpoints:   make(map[string]*models.Endpoint),
		WorkerPool:  make(chan chan string, params.MaxWorker),
		Queue:       make(chan string, params.MaxQueue),
		MaxWorker:   params.MaxWorker,
		MaxQueue:    params.MaxQueue,
		Signals:     make(chan os.Signal, 2),
		Quit:        make(chan bool),
		DB:          db,
		Logger:      logentry,
	}

	return dispatcher, nil
}

// stop signals programmatically
func (dispatcher *Dispatcher) Stop() {
	go func() {
		dispatcher.Quit <- true
	}()
}

// stop signals from system
func (dispatcher *Dispatcher) Signal() {
	go func() {

		// link system signal to the dispatcher signal
		signal.Notify(dispatcher.Signals, os.Interrupt, syscall.SIGTERM)

		// when receive syscall signal
		<-dispatcher.Signals

		// do a stopper
		dispatcher.Stop()
	}()
}

///////////////////////////////////////////////////////////////////////////////

package machine

import (
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// interface to implements to correctly works with this executor
type Machine interface {
	// permit to retrieve logger handler
	GetLogger() (*logrus.Logger, error)
	// permit to retrieve database handler
	GetDatabase(*logrus.Entry) (*pg.DB, error)
}

///////////////////////////////////////////////////////////////////////////////

// default machine definition
type DefaultMachine struct{}

// permit to retrieve logger handler
func (machine DefaultMachine) GetLogger() (*logrus.Logger, error) {

	// create a default logger
	var logger = logrus.New()

	// define the logger default level
	logger.Level = logrus.DebugLevel

	// define the logger default output
	logger.Out = os.Stdout

	return logger, nil
}

// permit to retrieve database handler
func (machine DefaultMachine) GetDatabase(logger *logrus.Entry) (*pg.DB, error) {

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
		logger.Println("Problem while check database connection")
		return nil, err
	}

	// build the query logger
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		// XXX : maybe only use UnformattedQuery option ( a debug flag ? )
		query, err := event.FormattedQuery()
		if err != nil {
			logger.Panic(err)
		}

		logger.WithField("sql_duration_ms", time.Since(event.StartTime)).Info(query)
	})

	return db, nil
}

///////////////////////////////////////////////////////////////////////////////

// to launch the machine executor with a default struct
func RunDefault(function string) {
	Run(DefaultMachine{}, function)
}

// to launch the machine executor
func Run(machine Machine, function string) {

	// create a dispatcher
	dispatcher, err := NewDispatcher(&DispatcherParams{
		Machine:   machine,
		Function:  function,
		MaxWorker: 20,
		MaxQueue:  5,
	})

	// deferred the database connection closes
	defer dispatcher.DB.Close()

	// dispatcher creation catch error
	if err != nil {
		logrus.Fatal(err)
	}

	// dispatcher logger
	log := dispatcher.Logger

	// get robot configuration
	err = dispatcher.GetRobotConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	// listen the channel
	go dispatcher.Signal()

	// launch a first task ID listing
	go dispatcher.GetTasks()

	//	// launch the database NOTIFY listener
	//	go dispatcher.Listen()

	// starting n number of workers
	for i := int64(0); i < dispatcher.MaxWorker; i++ {

		// create a new worker
		worker := NewWorker(i, dispatcher)

		// start it
		worker.Start()
	}

	// launch the dispatch
	log.Info("Worker dispatch started...")

	for {
		select {
            case ID := <-dispatcher.Queue:

                log.WithField(LABEL_TASK_ID, ID).Info("Dispatch task")

                // try to obtain a worker task channel that is available.
                // this will block until a worker is idle
                taskChannel := <-dispatcher.WorkerPool

                // dispatch the task to the worker task channel
                taskChannel <- ID

            case <-dispatcher.Quit:

                // we have received a signal to stop
                log.Info("Dispatch is stopping")

                // XXX : how to stop workers correctly

                return
		}
	}
}

// Label for logger fields
var (
    LABEL_WORKER_ID string = "worker_id"
    LABEL_TASK_ID string = "task_id"
    LABEL_STEP string = "step"
    LABEL_FUNCTION string = "function"
)


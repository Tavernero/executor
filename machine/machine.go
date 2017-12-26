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

	// XXX: get robot configuration

	// listen the channel
	go dispatcher.Signal()

	// XXX: launch a first task ID listing

	// XXX: launch the database NOTIFY listener

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
	LABEL_TASK_ID   string = "task_id"
	LABEL_STEP      string = "step"
	LABEL_FUNCTION  string = "function"
)

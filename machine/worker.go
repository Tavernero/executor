package machine

import (
	"github.com/sirupsen/logrus"
)

///////////////////////////////////////////////////////////////////////////////

// the worker executes the task process
type Worker struct {
	// ID of the local worker
	ID int64
	// Pool of worker to notify our free time
	Pool chan chan string
	// Channel of the task IDs
	Channel chan string
	// Parent dispatcher
	Dispatcher *Dispatcher
	// Quit chan
	Quit chan bool
}

// worker creation handler
func NewWorker(ID int64, dispatcher *Dispatcher) Worker {
	return Worker{
		ID:         ID,
		Pool:       dispatcher.WorkerPool,
		Channel:    make(chan string),
		Dispatcher: dispatcher,
		Quit:       make(chan bool),
	}
}

// start method starts the run loop for the worker
// listening for a quit channel in case we need to stop it
func (worker *Worker) Start() {
	go func() {

		// infinite loop
		for {

			// register the current worker into the worker queue
			worker.Pool <- worker.Channel

			// read from the channel
			select {

			case ID := <-worker.Channel:

				logrus.Info("Worker works on task")

                // XXX: do the action on hhttp calls

			case <-worker.Quit:

				logrus.Info("Worker quits")

				// we have received a signal to stop
				// exit this function
				return
			}
		}
	}()
}

// stop signals the worker to stop listening for work requests.
func (worker *Worker) Stop() {
	go func() {
		worker.Quit <- true
	}()
}

///////////////////////////////////////////////////////////////////////////////

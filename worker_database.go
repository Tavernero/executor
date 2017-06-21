package main

import (
    "gopkg.in/gorp.v2"
)

// ================================================= //
// ================================================= //

// Worker represents the worker that executes the task
type DatabaseWorker struct {
    DatabaseWorkerPool  chan chan *gorp.DbMap
    DatabaseChannel     chan Task
    quit                chan bool
    connector           *gorp.DbMap
}

// Create a new worker
func NewDatabaseWorker(databaseWorkerPool chan chan *gorp.DbMap) DatabaseWorker {
    return DatabaseWorker{
        DatabaseWorkerPool: databaseWorkerPool,
        DatabaseChannel: make(chan Task),
        quit: make(chan bool) }
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w DatabaseWorker) Start() {
    // Retrieve a new connection for each worker instance
    w.connector = initDb();
}

// Stop signals the worker to stop listening for work requests.
func (w DatabaseWorker) Stop() {
    go func() {
        w.quit <- true
    }()
}

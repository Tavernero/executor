package main

import (
    "gopkg.in/gorp.v2"
)

// ================================================= //
// ================================================= //

// Worker represents the worker that executes the task
type DatabaseWorker struct {
    DatabaseWorkerPool  chan chan *gorp.DbMap
    DatabaseChannel chan Task
    quit        chan bool
}

// Create a new worker
func NewDatabaseWorker(databaseWorkerPool chan chan *gorp.DbMap) DatabaseWorker {
    return DatabaseWorker{
        DatabaseWorkerPool: databaseWorkerPool,
        DatabaseChannel: make(chan Task),
        quit: make(chan bool)}
}

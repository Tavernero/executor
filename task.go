package main

import (
    "fmt"
    "strconv"
    "math/rand"
    "time"
)

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

// ================================================= //
// ================================================= //

// function called when we launch a task object
func (t *Task) WorkingOn() error {
    randomer := rand.Intn(2000)
    time.Sleep(time.Duration(randomer) * time.Millisecond)
    fmt.Println("work done in " + strconv.Itoa( randomer ) + " ms for '" + strconv.Itoa( t.ID ) + "'" )
    return nil
}

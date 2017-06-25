package main

import (

)

// Task to work
type Task struct {
    ID int `db:"id, primarykey" json:"id"`
    Function string `db:"function" json:"function"`
    Name string `db:"name" json:"name"`
    Step string `db:"step" json:"step"`
    Status string `db:"status" json:"status"`
    Retry int `db:"retry" json:"retry"`
    Arguments JsonB `db:"arguments" json:"arguments"`
    Buffer JsonB `db:"buffer" json:"buffer"`
}

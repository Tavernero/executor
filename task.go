package main

import (
    "time"
//    "gopkg.in/gorp.v2"
)

// Task to work
type Task struct {
    ID int `db:"id, primarykey" json:"id"`
    Function string `db:"function" json:"function"`
    Name string `db:"name" json:"name"`
    Step string `db:"step" json:"step"`
    Status string `db:"status" json:"status"`
    Retry int `db:"retry" json:"retry"`
    Comment string `db:"comment" json:"comment"`
//    CreationDate *time.Time `db:"creation_date" json:"creation_date"`
//    TodoDate *time.Time `db:"todo_date" json:"todo_date"`


//    LastUpdate *time.Time `db:"last_update" json:"last_update"`


    TodoDate *time.Time `db:"todo_date" json:"todo_date"`
    LastUpdate *time.Time `db:"last_update" json:"last_update"`
    DoneDate *time.Time `db:"done_date" json:"done_date"`



//    DoneDate *time.Time `db:"done_date" json:"done_date"`
    Arguments JsonB `db:"arguments" json:"arguments"`
    Buffer JsonB `db:"buffer" json:"buffer"`
}

package main

import (
    "time"
)

type Definition struct {
    ID int64 `db:"id, primarykey, autoincrement" json:"id"`
    Function string `db:"function" json:"function"`
    Status string `db:"status" json:"status"`
    Sequence Sequence `db:"sequence" json:"sequence"`
    CreationDate *time.Time `db:"creation_date" json:"creation_date"`
    LastUpdate *time.Time `db:"last_update" json:"last_update"`
}

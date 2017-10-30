package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Robot is the go counterpart of table robot
type Robot struct {
	TableName struct{} `sql:"robot"`

	ID           string                 `sql:"id"`
	Function     string                 `sql:"function,notnull"`
	Version      int64                  `sql:"version,notnull"`
	Status       bool                   `sql:"status,notnull"`
	Definition   map[string]interface{} `sql:"definition,notnull"`
	CreationDate time.Time              `sql:"creation_date,notnull"`
	LastUpdate   time.Time              `sql:"last_update,notnull"`
}

const (
	ColRobot_ID           = `"id"`
	ColRobot_Function     = `"function"`
	ColRobot_Version      = `"version"`
	ColRobot_Status       = `"status"`
	ColRobot_Definition   = `"definition"`
	ColRobot_CreationDate = `"creation_date"`
	ColRobot_LastUpdate   = `"last_update"`
)

const (
	TblRobot_ID           = `"robot"."id"`
	TblRobot_Function     = `"robot"."function"`
	TblRobot_Version      = `"robot"."version"`
	TblRobot_Status       = `"robot"."status"`
	TblRobot_Definition   = `"robot"."definition"`
	TblRobot_CreationDate = `"robot"."creation_date"`
	TblRobot_LastUpdate   = `"robot"."last_update"`
)

///////////////////////////////////////////////////////////////////////////////

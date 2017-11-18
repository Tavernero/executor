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

	ID           string     `sql:"id"`
	Function     string     `sql:"function,notnull"`
	Version      int64      `sql:"version,notnull"`
	Status       bool       `sql:"status,notnull"`
	Definition   Definition `sql:"definition,notnull"`
	CreationDate time.Time  `sql:"creation_date,notnull"`
	LastUpdate   time.Time  `sql:"last_update,notnull"`
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

type Step struct {
	Name       string `json:"name"`
	EndpointID string `json:"endpoint_id"`
}

///////////////////////////////////////////////////////////////////////////////

type Definition struct {
	Sequence []Step `json:"sequence"`
}

var (
	_ sql.Scanner   = &Definition{}
	_ driver.Valuer = Definition{}
)

// Scan implements the sql.Scanner interface.
func (e *Definition) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case Definition:
		// Need to deep copy struct
		*e = Definition{
			Sequence: make([]Step, len(v.Sequence)),
		}
		for idx, step := range v.Sequence {
			e.Sequence[idx] = step
		}
	case []byte:
		*e = Definition{} // json.Unmarshal does not reset the struct content
		err = json.Unmarshal(v, e)
	case string:
		*e = Definition{} // json.Unmarshal does not reset the struct content
		err = json.Unmarshal([]byte(v), e)
	default:
		return fmt.Errorf("Can't convert %T to models.Definition", value)
	}
	return
}

var emptySlice = []Step{}

// Value implements the driver driver.Valuer interface.
func (e Definition) Value() (driver.Value, error) {
	if e.Sequence == nil {
		e.Sequence = emptySlice
	}
	jsonBytes, _ := json.Marshal(e)
	return string(jsonBytes), nil
}

///////////////////////////////////////////////////////////////////////////////

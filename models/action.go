package models

import (
	"database/sql/driver"
	"fmt"
)

// Action type is the go counterpart of enum for action key of endpoint response
type Action string

// Scan implements the sql.Scanner interface.
func (e *Action) Scan(value interface{}) error {
	if value == nil {
		*e = Action_NULL
		return nil
	}

	switch v := value.(type) {
	case string:
		*e = Action(v)
		return nil
	case []byte:
		*e = Action(v)
		return nil
	}
	return fmt.Errorf("Can't convert %T to Action", value)
}

// Value implements the driver driver.Valuer interface.
func (e Action) Value() (driver.Value, error) {
	if e == Action_NULL {
		return nil, nil
	}
	return string(e), nil
}

const (
	Action_NULL       Action = ""
	Action_DONE       Action = "DONE"
	Action_GOTO       Action = "GOTO"
	Action_NEXT       Action = "NEXT"
	Action_GOTO_LATER Action = "GOTO_LATER"
	Action_NEXT_LATER Action = "NEXT_LATER"
	Action_RETRY      Action = "RETRY"
	Action_RETRY_NOW  Action = "RETRY_NOW"
	Action_ERROR      Action = "ERROR"
	Action_PROBLEM    Action = "PROBLEM"
	Action_CANCELED   Action = "CANCELED"
)

///////////////////////////////////////////////////////////////////////////////

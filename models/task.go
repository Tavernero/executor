package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Task is the go counterpart of table task
type Task struct {
	TableName struct{} `sql:"task"`

	ID           string     `sql:"id"`
	Version      int64      `sql:"version,notnull"`
	Context      string     `sql:"context,notnull"`
	Function     string     `sql:"function,notnull"`
	Step         string     `sql:"step,notnull"`
	Status       TaskStatus `sql:"status,notnull"`
	Retry        int64      `sql:"retry,notnull"`
	CreationDate time.Time  `sql:"creation_date,notnull"`
	LastUpdate   time.Time  `sql:"last_update,notnull"`
	TodoDate     time.Time  `sql:"todo_date,notnull"`
	DoneDate     *time.Time `sql:"done_date"`
	Arguments    JsonB      `sql:"arguments,notnull"`
	Buffer       JsonB      `sql:"buffer,notnull"`
	Comment      string     `sql:"comment"`
}

const (
	ColTask_ID           = `id`
	ColTask_Version      = `version`
	ColTask_Context      = `context`
	ColTask_Function     = `function`
	ColTask_Step         = `step`
	ColTask_Status       = `status`
	ColTask_Retry        = `retry`
	ColTask_CreationDate = `creation_date`
	ColTask_LastUpdate   = `last_update`
	ColTask_TodoDate     = `todo_date`
	ColTask_DoneDate     = `done_date`
	ColTask_Arguments    = `arguments`
	ColTask_Buffer       = `buffer`
	ColTask_Comment      = `comment`
)

const (
	TblTask_ID           = `"task"."id"`
	TblTask_Version      = `"task"."version"`
	TblTask_Context      = `"task"."context"`
	TblTask_Function     = `"task"."function"`
	TblTask_Step         = `"task"."step"`
	TblTask_Status       = `"task"."status"`
	TblTask_Retry        = `"task"."retry"`
	TblTask_CreationDate = `"task"."creation_date"`
	TblTask_LastUpdate   = `"task"."last_update"`
	TblTask_TodoDate     = `"task"."todo_date"`
	TblTask_DoneDate     = `"task"."done_date"`
	TblTask_Arguments    = `"task"."arguments"`
	TblTask_Buffer       = `"task"."buffer"`
	TblTask_Comment      = `"task"."comment"`
)

///////////////////////////////////////////////////////////////////////////////

// TaskStatus type is the go counterpart of DB enum task_status
type TaskStatus string

// Scan implements the sql.Scanner interface.
func (e *TaskStatus) Scan(value interface{}) error {
	if value == nil {
		*e = TaskStatus_NULL
		return nil
	}

	switch v := value.(type) {
	case string:
		*e = TaskStatus(v)
		return nil
	case []byte:
		*e = TaskStatus(v)
		return nil
	}
	return fmt.Errorf("Can't convert %T to TaskStatus", value)
}

// Value implements the driver driver.Valuer interface.
func (e TaskStatus) Value() (driver.Value, error) {
	if e == TaskStatus_NULL {
		return nil, nil
	}
	return string(e), nil
}

const (
	TaskStatus_MISTAKE  TaskStatus = "MISTAKE"
	TaskStatus_FAULT    TaskStatus = "FAULT"
	TaskStatus_NULL     TaskStatus = ""
	TaskStatus_TODO     TaskStatus = "TODO"
	TaskStatus_DOING    TaskStatus = "DOING"
	TaskStatus_ERROR    TaskStatus = "ERROR"
	TaskStatus_PROBLEM  TaskStatus = "PROBLEM"
	TaskStatus_DONE     TaskStatus = "DONE"
	TaskStatus_CANCELED TaskStatus = "CANCELED"
)

///////////////////////////////////////////////////////////////////////////////

package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PropertyMap for catch JSONB from databases
type JsonB map[string]interface{}

func (p JsonB) Value() (driver.Value, error) {

	j, err := json.Marshal(p)

	return string(j), err
}

func (p *JsonB) Scan(src interface{}) error {

	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}

	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

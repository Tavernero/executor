package main

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
)

// for catch JSONB from databases
type Sequence []Step

func (p Sequence) Value() (driver.Value, error) {

    j, err := json.Marshal(p)

    return j, err
}

func (p *Sequence) Scan(src interface{}) error {

    source, ok := src.([]byte)

    if !ok {
        return errors.New("Type assertion .([]byte) failed.")
    }

    err := json.Unmarshal(source, &p)

    if err != nil {
        return err
    }

    return nil
}

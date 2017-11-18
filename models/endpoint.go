package models

import (
	"time"
)

// Endpoint is the go counterpart of table endpoint
type Endpoint struct {
	TableName struct{} `sql:"endpoint"`

	ID           string    `sql:"id"`
	Version      int64     `sql:"version,notnull"`
	Name         string    `sql:"name,notnull"`
	Method       string    `sql:"method,notnull"`
	URL          string    `sql:"url,notnull"`
	CreationDate time.Time `sql:"creation_date,notnull"`
	LastUpdate   time.Time `sql:"last_update,notnull"`
}

const (
	ColEndpoint_ID           = `id`
	ColEndpoint_Version      = `version`
	ColEndpoint_Name         = `name`
	ColEndpoint_Method       = `method`
	ColEndpoint_URL          = `url`
	ColEndpoint_CreationDate = `creation_date`
	ColEndpoint_LastUpdate   = `last_update`
)

const (
	TblEndpoint_ID           = `"endpoint"."id"`
	TblEndpoint_Version      = `"endpoint"."version"`
	TblEndpoint_Name         = `"endpoint"."name"`
	TblEndpoint_Method       = `"endpoint"."method"`
	TblEndpoint_URL          = `"endpoint"."url"`
	TblEndpoint_CreationDate = `"endpoint"."creation_date"`
	TblEndpoint_LastUpdate   = `"endpoint"."last_update"`
)

///////////////////////////////////////////////////////////////////////////////

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameGooseDbVersion = "goose_db_version"

// GooseDbVersion mapped from table <goose_db_version>
type GooseDbVersion struct {
	ID        int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	VersionID int64     `gorm:"column:version_id;type:bigint;not null" json:"version_id"`
	IsApplied bool      `gorm:"column:is_applied;type:boolean;not null" json:"is_applied"`
	Tstamp    time.Time `gorm:"column:tstamp;type:timestamp without time zone;not null;default:now()" json:"tstamp"`
}

// TableName GooseDbVersion's table name
func (*GooseDbVersion) TableName() string {
	return TableNameGooseDbVersion
}

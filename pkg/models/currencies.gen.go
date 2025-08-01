// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCurrency = "currencies"

// Currency mapped from table <currencies>
type Currency struct {
	UUID           string         `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid"`
	CountryName    string         `gorm:"column:country_name;type:character varying(255);not null;index:idx_currencies_country_name,priority:1" json:"country_name"`
	CurrencyName   string         `gorm:"column:currency_name;type:character varying(255);not null;index:idx_currencies_currency_name,priority:1" json:"currency_name"`
	CountryFlag    string         `gorm:"column:country_flag;type:character varying(10);not null" json:"country_flag"`
	CurrencySymbol string         `gorm:"column:currency_symbol;type:character varying(10);not null" json:"currency_symbol"`
	CurrencyCode   string         `gorm:"column:currency_code;type:character varying(3);not null;uniqueIndex:idx_currencies_currency_code,priority:1" json:"currency_code"`
	CreatedAt      *time.Time     `gorm:"column:created_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at;type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp with time zone;index:idx_currencies_deleted_at,priority:1" json:"deleted_at"`
}

// TableName Currency's table name
func (*Currency) TableName() string {
	return TableNameCurrency
}

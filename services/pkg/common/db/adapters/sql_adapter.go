package adapters

import (
	"gorm.io/gorm"
)

type SqlDataAdapter interface {
	GetDB() (*gorm.DB, error)
	Migrate(...interface{}) error
	Ping() error
}

package model

import (
	"time"

	"gorm.io/datatypes"
)

type School struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"index"`
	Address string
}

type Log struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"index"`
	Action    string
	Data      datatypes.JSON
}

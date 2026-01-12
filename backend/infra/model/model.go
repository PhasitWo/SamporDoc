package model

import (
	"time"

	"gorm.io/datatypes"
)

type Customer struct {
	ID      uint64  `gorm:"primaryKey;autoIncrement"`
	Name    string  `gorm:"index" json:"name"`
	Address *string `json:"address" ts_type:"string | null"`
}

type Shop struct {
	ID                 uint64  `gorm:"primaryKey;autoIncrement"`
	Slug               string  `gorm:"uniqueIndex" json:"slug"`
	Name               string  `json:"name"`
	SortingLevel       uint64  `json:"sortingLevel"`
	ReceiptFormPath    *string `json:"receiptFormPath" ts_type:"string | null"`
	ReceiptControlPath *string `json:"receiptControlPath" ts_type:"string | null"`
}

type Log struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time `gorm:"index"`
	Status        string
	Action        string
	Data          datatypes.JSON
	CorrelationID int64 `gorm:"index"`
}

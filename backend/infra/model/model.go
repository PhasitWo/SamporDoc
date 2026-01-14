package model

import (
	"time"

	"gorm.io/datatypes"
)

type Customer struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement"`
	Name            string  `gorm:"index" json:"name"`
	Address         *string `json:"address" ts_type:"string | null"`
	HeadCheckerName *string `json:"headCheckerName" ts_type:"string | null"`
	Checker1Name    *string `json:"checker1Name" ts_type:"string | null"`
	Checker2Name    *string `json:"checker2Name" ts_type:"string | null"`
	ObjectName      *string `json:"objectName" ts_type:"string | null"`
	HeadObjectName  *string `json:"headObjectName" ts_type:"string | null"`
	BossName        *string `json:"bossName" ts_type:"string | null"`
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

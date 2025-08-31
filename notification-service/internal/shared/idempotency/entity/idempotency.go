package entity

import "gorm.io/gorm"

type Idempotency struct {
	gorm.Model
	Key      string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Request  string `gorm:"type:text;not null"`
	Response string `gorm:"type:text"`
}

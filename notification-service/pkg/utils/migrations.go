package utils

import (
	"gorm.io/gorm"
)

func InitMigrations(db *gorm.DB) {
	db.AutoMigrate()
}

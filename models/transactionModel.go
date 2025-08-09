package models

import "gorm.io/gorm"

type Transaction struct {
    gorm.Model
    UserID     uint     `gorm:"not null;index"`
    CategoryID *uint    `gorm:"index"`
    Amount     float64  `gorm:"not null"`
    Description string  `gorm:"type:varchar(255)"`
}

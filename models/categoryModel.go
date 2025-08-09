package models

import "gorm.io/gorm"

type Category struct {
    gorm.Model
    UserID       uint         `gorm:"not null;index"`
    Name         string       `gorm:"type:varchar(100);not null"`
    Transactions []Transaction `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
}

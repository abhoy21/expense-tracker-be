package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
    Email    string `gorm:"not null;uniqueIndex"`
    Password string `gorm:"not null"`
    Categories []Category  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Transactions []Transaction `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

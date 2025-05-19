package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `json:"id" gorm:"unique;not null"`
	Email    string  `json:"email" gorm:"unique;not null"`
	Password string  `json:"-" gorm:"not null"`
	Role     *string `json:"role"`
}

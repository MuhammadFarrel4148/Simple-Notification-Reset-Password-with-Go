package models

import "gorm.io/gorm"

type Gouser struct {
	gorm.Model
	Email		string 		`gorm:"unique;not null"`
	Password	string 		`gorm:"not null"`
}

type Gootp struct {
	Id			uint		`gorm:"primaryKey;autoIncrement"`
	Email		string		`gorm:"not null"`
	Otp			string		`gorm:"not null"`
}
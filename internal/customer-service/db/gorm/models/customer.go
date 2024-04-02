package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	*gorm.Model
	Firstname   string
	Middlename  *string
	Lastname    string
	DateOfBirth time.Time
}

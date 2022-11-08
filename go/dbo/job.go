package dbo

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Data string
}

package dbo

import (
	"gorm.io/gorm"
)

type JobType int32

const (
	Rectangle JobType = iota
	Stator
)

type JobStatus int32

const (
	New JobStatus = iota
	Meshing
)

type Job struct {
	gorm.Model

	Type   JobType   `gorm:"default:0"`
	Status JobStatus `gorm:"default:0"`
	Data   string
}

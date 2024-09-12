package domain

import (
	"time"

	"gorm.io/gorm"
)

type ProjectDocument struct {
	Id           int64          `gorm:"primaryKey,autoIncrement"`
	ProjectId    int64          `gorm:"index"` //
	Name         string         ``
	Url          string         ``
	DocumentType string         ``
	CreatedAt    time.Time      `gorm:"autoCreateTime:true"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime:true"`
	DeletedAt    gorm.DeletedAt ``
} //@name ProjectDescription

func (*ProjectDocument) TableName() string { return TableNameProjectDocument }

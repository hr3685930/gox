package queue

import "github.com/golang-module/carbon"

type FailedJobs struct {
	Id         int64         `gorm:"column:id" json:"id"`
	Connection string        `gorm:"column:connection;types:text"`
	Topic      string        `gorm:"column:topic;types:text"`
	Queue      string        `gorm:"column:queue;types:text"`
	Message    string        `gorm:"column:message;types:text"`
	Exception  string        `gorm:"column:exception;types:longText"`
	Stack      string        `gorm:"column:stack;types:longText"`
	FiledAt    carbon.Carbon `gorm:"column:failed_at"`
}

func (*FailedJobs) TableName() string {
	return "failed_jobs"
}

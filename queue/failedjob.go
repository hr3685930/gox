package queue

import "github.com/golang-module/carbon"

type FailedJobs struct {
	Id         int64         `gorm:"column:id" json:"id"`
	Connection string        `gorm:"column:connection;types:text" json:"connection"`
	Topic      string        `gorm:"column:topic;types:text" json:"topic"`
	Queue      string        `gorm:"column:queue;types:text" json:"queue"`
	Message    string        `gorm:"column:message;types:text" json:"message"`
	Exception  string        `gorm:"column:exception;types:longText" json:"exception"`
	Stack      string        `gorm:"column:stack;types:longText" json:"stack"`
	FiledAt    carbon.Carbon `gorm:"column:failed_at" json:"filed_at"`
}

func (*FailedJobs) TableName() string {
	return "failed_jobs"
}

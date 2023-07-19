package models

import (
	"gorm.io/gorm"
	"time"
)

// 填写面试表单信息
type InterviewRecord struct {
	gorm.Model
	InterviewUsername  string    `json:"interview_username" gorm:"not null;comment:面试人员"`
	InterviewCompany   string    `json:"interview_company" gorm:"not null;comment:面试公司"`
	InterviewPosition  string    `json:"interview_position" gorm:"not null;comment:面试岗位"`
	InterviewLocation  string    `json:"interview_location" gorm:"not null;comment:面试地点"`
	InterviewType      int       `json:"interview_type" gorm:"not null;comment:面试类型"`
	InterviewApproach  int       `json:"interview_approach" gorm:"not null;comment:面试方式"`
	InterviewTime      time.Time `json:"interview_time" gorm:"not null;comment:面试时间"`
	InterviewConfusion string    `json:"interview_confusion" gorm:"not null;comment:面试困惑"`
}

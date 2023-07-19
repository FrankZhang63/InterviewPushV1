package request

import "time"

// 展示的数据
type InterviewMsg struct {
	InterviewUsername  string    `json:"interview_username" validate:"required" form:"interview_username"`
	InterviewCompany   string    `json:"interview_company" validate:"required" form:"interview_company"`
	InterviewPosition  string    `json:"interview_position" validate:"required" form:"interview_position"`
	InterviewLocation  string    `json:"interview_location" validate:"required" form:"interview_location"`
	InterviewType      int       `json:"interview_type" validate:"required" form:"interview_type"`
	InterviewApproach  int       `json:"interview_approach" validate:"required" form:"interview_approach"`
	InterviewTime      time.Time `json:"interview_time" validate:"required" form:"interview_time"`
	InterviewConfusion string    `json:"interview_confusion" validate:"required" form:"interview_confusion"`
}

// 分页展示
type InterviewPageSize struct {
	InterviewSelect
	Page int `json:"page" form:"page"` //页码
}

// 多条件查询
type InterviewSelect struct {
	InterviewUsername  string    `json:"interview_username" form:"interview_username"`
	InterviewCompany   string    `json:"interview_company" form:"interview_company"`
	InterviewPosition  string    `json:"interview_position" form:"interview_position"`
	InterviewLocation  string    `json:"interview_location" form:"interview_location"`
	InterviewType      int       `json:"interview_type" form:"interview_type"`
	InterviewApproach  int       `json:"interview_approach" form:"interview_approach"`
	InterviewBeginTime time.Time `json:"interview_begin_time" form:"interview_begin_time"`
	InterviewEndTime   time.Time `json:"interview_end_time" form:"interview_end_time"`
}

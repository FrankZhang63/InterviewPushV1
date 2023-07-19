package mysql

import (
	"InterviewPush/models"
	"InterviewPush/models/request"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// CreateMessage 添加一条面试记录
func CreateMessage(msg request.InterviewMsg) (err error) {
	result := db.Where("interview_username = ? AND interview_time = ?", msg.InterviewUsername, msg.InterviewTime).FirstOrCreate(&models.InterviewRecord{
		Model:              gorm.Model{},
		InterviewUsername:  msg.InterviewUsername,
		InterviewCompany:   msg.InterviewCompany,
		InterviewPosition:  msg.InterviewPosition,
		InterviewLocation:  msg.InterviewLocation,
		InterviewType:      msg.InterviewType,
		InterviewApproach:  msg.InterviewApproach,
		InterviewTime:      msg.InterviewTime,
		InterviewConfusion: msg.InterviewConfusion,
	})
	if result.RowsAffected == 0 {
		zap.L().Error("已经存在", zap.Error(err))
		return ErrorMsgExit
	}
	if result.Error != nil {
		zap.L().Error("CreateMessage failed:", zap.Error(err))
		return ErrorInsertFailed
	}
	return
}

// SelectMessage 查询面试记录分页 多条件
func SelectMessage(info request.InterviewPageSize) (data []request.InterviewMsg, total int64, err error) {
	//limit := info.PageSize
	offset := 10 * (info.Page - 1)
	db := db.Model(&models.InterviewRecord{})
	var interviewMsg []request.InterviewMsg
	// 面试名字
	if info.InterviewUsername != "" {
		db = db.Where("interview_username LIKE ?", "%"+info.InterviewUsername+"%")
	}
	// 面试公司
	if info.InterviewCompany != "" {
		db = db.Where("interview_company LIKE ?", "%"+info.InterviewCompany+"%")
	}
	// 面试岗位
	if info.InterviewPosition != "" {
		fmt.Println(info.InterviewPosition)
		db = db.Where("interview_position LIKE ?", "%"+info.InterviewPosition+"%")
	}
	// 面试地点
	if info.InterviewLocation != "" {
		db = db.Where("interview_location LIKE ?", "%"+info.InterviewLocation+"%")
	}
	// 面试类型
	if info.InterviewType != 0 {
		db = db.Where("interview_type = ?", info.InterviewType)
	}
	// 面试方式
	if info.InterviewApproach != 0 {
		db = db.Where("interview_approach = ?", info.InterviewApproach)
	}
	// 面试时间段
	// 面试时间都填了
	if !info.InterviewEndTime.IsZero() && !info.InterviewBeginTime.IsZero() {
		db = db.Where("interview_time >= ? AND interview_time <= ?", info.InterviewBeginTime, info.InterviewEndTime)
	} else if info.InterviewEndTime.IsZero() { //结束时间没填
		db = db.Where("interview_time >= ?", info.InterviewBeginTime)
	} else if info.InterviewBeginTime.IsZero() { //开始时间没填
		db = db.Where("interview_time <= ?", info.InterviewEndTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		zap.L().Error("SelectAllMsg sql failed:", zap.Error(err))
		return
	}
	err = db.Debug().Limit(10).Offset(offset).Order("interview_time").Scan(&interviewMsg).Error
	return interviewMsg, total, err
}

// 根据昵称查询此人的面试记录
func SelectMessageByName(Name string) (data []request.InterviewMsg, err error) {
	err = db.Table("interview_records").Where("interview_username = ? AND interview_time >= ?", Name, time.Now()).Scan(&data).Error
	if err != nil {
		zap.L().Error("SelectMessageByName sql failed:", zap.Error(err))
		return
	}
	return data, err
}

package logic

import (
	"InterviewPush/dao/mysql"
	"InterviewPush/models/request"
	"go.uber.org/zap"
)

// AddMessage 收集信息
func AddMessage(msg request.InterviewMsg) (err error) {
	err = mysql.CreateMessage(msg)
	if err != nil {
		zap.L().Error("CreateMsg logic failed:", zap.Error(err))
		return
	}
	return
}

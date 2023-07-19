package controller

import (
	"InterviewPush/logic"
	reqm "InterviewPush/models/request"
	"InterviewPush/models/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AddMessageHandler 填写表单收集信息
func AddMessageHandler(c *gin.Context) {
	// 参数绑定
	p := new(reqm.InterviewMsg)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("create comment with invalid param", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.CodeInvalidParams, "参数绑定错误")
		return
	}
	// 参数验证
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		zap.L().Error("有字段为空", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.CodeInvalidParams, "存在空参数")
		return
	}
	fmt.Println("用户名:", p.InterviewUsername, "公司:", p.InterviewCompany, "岗位:", p.InterviewPosition, "面试地点:", p.InterviewLocation, "类型:", p.InterviewType, "类型:", p.InterviewApproach, "面试时间:", p.InterviewTime)
	// service操作
	err := logic.AddMessage(*p)
	if err != nil {
		zap.L().Error("addmessage failed:", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	PushInformation(p) //调用机器人发送信息
	response.ResponseSuccess(c, "添加记录成功")
}

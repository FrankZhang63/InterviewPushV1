package Interview_admin

import (
	"InterviewPush/models/request"
	"InterviewPush/models/response"
	"InterviewPush/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 登录
func LoginHandler(c *gin.Context) {
	var p request.AdminLogin
	err := c.ShouldBindJSON(&p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	// 参数验证
	validate := validator.New()
	if err := validate.Struct(p); err != nil {
		zap.L().Error("有字段为空", zap.Error(err))
		response.ResponseErrorWithMsg(c, response.CodeInvalidParams, "用户名或密码为空")
		return
	}
	if p.Username == "wtj" && p.Password == "123456" {
		token, _, _ := jwt.GenToken(p.Username)
		fmt.Println("登录成功", p.Username, "+++++", p.Password)
		response.ResponseSuccess(c, token)
		return
	}
	fmt.Println("登录失败", p.Username, "+++++", p.Password)
	response.ResponseErrorWithMsg(c, response.CodeInvalidPassword, "用户名或密码错误")
}

// 登出
func LogoutHandler(c *gin.Context) {
	response.ResponseSuccess(c, "删除成功")
}

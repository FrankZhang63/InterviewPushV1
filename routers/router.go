package routers

import (
	"InterviewPush/controller"
	"InterviewPush/controller/Interview_admin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(controller.Cors())
	v1 := r.Group("/interview")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "I'll love you till8 I die")
	})
	//前台收集表单自动登录
	v1.POST("/add", controller.AddMessageHandler)
	//后台管理员登录
	v1.POST("/login", Interview_admin.LoginHandler)
	//接收信息
	r.POST("/reception", controller.ReceptionHandler)
	//登录验证token
	v1.Use(controller.JWTAuthMiddleware())
	{
		//全查
		v1.GET("/showmsg", Interview_admin.ShowMsgHandler)
		//登出
		v1.GET("/logout", Interview_admin.LogoutHandler)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "独一无二的错误",
		})
	})
	return r
}

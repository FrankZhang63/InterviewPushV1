package Interview_admin

import (
	"InterviewPush/dao/mysql"
	"InterviewPush/models/request"
	"InterviewPush/models/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ShowMsgHandler 查询展示列表
func ShowMsgHandler(c *gin.Context) {
	var pageInfo request.InterviewPageSize
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		zap.L().Error("参数绑定错误", zap.Error(err))
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	if list, total, err := mysql.SelectMessage(pageInfo); err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	} else {
		response.ResponseSuccess(c, response.PageResult{
			List:  list,
			Total: total,
			Page:  pageInfo.Page,
		})
	}
}

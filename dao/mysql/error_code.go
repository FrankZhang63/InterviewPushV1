package mysql

import "errors"

var (
	ErrorMsgExit       = errors.New("面试信息已存在")
	ErrorUserExit      = errors.New("用户名已存在")
	ErrorUserNotExit   = errors.New("用户不存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorInvalidID     = errors.New("无效的ID")
	ErrorQueryFailed   = errors.New("查询数据失败")
	ErrorInsertFailed  = errors.New("添加面试记录失败")
)

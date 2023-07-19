package request

// 登录信息
type AdminLogin struct {
	Username string `json:"username" validate:"required" form:"username"`
	Password string `json:"password" validate:"required" form:"password"`
}

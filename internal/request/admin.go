package request

type InitPwdRequest struct {
	Password string `json:"password" form:"password" binding:"required"`
}

type AdminLoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

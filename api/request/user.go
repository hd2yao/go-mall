package request

type UserRegister struct {
	LoginName       string `json:"login_name" binding:"required,e164|email"` // 必须是手机号或邮箱
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"max=30"`
	Slogan          string `json:"slogan" binding:"max=30"`
	Avatar          string `json:"avatar" binding:"max=100"`
}

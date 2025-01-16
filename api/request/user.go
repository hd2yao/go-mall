package request

type UserRegister struct {
	LoginName       string `json:"login_name" binding:"required,e164|email"` // 必须是手机号或邮箱
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"max=30"`
	Slogan          string `json:"slogan" binding:"max=30"`
	Avatar          string `json:"avatar" binding:"max=100"`
}

type UserLogin struct {
	Body struct {
		LoginName string `json:"login_name" binding:"required,e164|email"`
		Password  string `json:"password" binding:"required,min=8"`
	}
	Header struct {
		Platform string `json:"platform" header:"platform" binding:"required,oneof=APP H5"`
	}
}

type UserInfoUpdate struct {
	Nickname string `json:"nickname" binding:"max=30"`
	Slogan   string `json:"slogan" binding:"max=30"`
	Avatar   string `json:"avatar" binding:"max=100"`
}

type PasswordResetApply struct {
	LoginName string `json:"login_name" binding:"required,e164|email"`
}

type PasswordReset struct {
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Token           string `json:"password_reset_token" binding:"required"`
	Code            string `json:"password_reset_code" binding:"required"`
}

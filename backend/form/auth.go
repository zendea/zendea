package form

// AuthSigninForm auth signin form
type AuthSigninForm struct {
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	CaptchaID   string  `form:"captchaId" json:"captchaId" binding: "required"`
	CaptchaCode string `form:"captchaCode" json:"captchaCode"`
}

// AuthSignupForm auth signup form
type AuthSignupForm struct {
	Email       string `form:"email" json:"email" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Nickname    string `form:"nickname" json:"nickname"`
	Password    string `form:"password" json:"password" binding:"required"`
	RePassword  string `form:"rePassword" json:"rePassword" binding:"required"`
	CaptchaID   string `form:"captchaId" json:"captchaId" binding: "required"`
	CaptchaCode string `form:"captchaCode" json:"captchaCode"`
}

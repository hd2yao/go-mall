package controller

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/appservice"
)

func RegisterUser(c *gin.Context) {
	userRequest := new(request.UserRegister)
	if err := c.ShouldBind(userRequest); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 密码复杂度验证
	if !util.PasswordComplexityVerify(userRequest.Password) {
		// Validator 验证通过后再应用 密码复杂度这样的特殊验证
		logger.New(c).Warn("RegisterUserError", "err", "密码复杂度不满足要求", "password", userRequest.Password)
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}

	// 注册用户
	userSvc := appservice.NewUserAppSvc(c)
	err := userSvc.UserRegister(userRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			app.NewResponse(c).Error(errcode.ErrUserNameOccupied)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
	return
}

func LoginUser(c *gin.Context) {
	loginRequest := new(request.UserLogin)
	if err := c.ShouldBindJSON(&loginRequest.Body); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if err := c.ShouldBindHeader(&loginRequest.Header); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	// 用户登录
	userSvc := appservice.NewUserAppSvc(c)
	token, err := userSvc.UserLogin(loginRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotRight) {
			app.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else if errors.Is(err, errcode.ErrUserInvalid) {
			app.NewResponse(c).Error(errcode.ErrUserInvalid)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		logger.New(c).Error("LoginUserError", "err", err)
		return
	}

	app.NewResponse(c).Success(token)
	return
}

func LogoutUser(c *gin.Context) {
	// 通过中间件从 token 中获取用户信息，并设置到 context 中
	userId := c.GetInt64("user_id")
	platform := c.GetString("platform")
	userSvc := appservice.NewUserAppSvc(c)
	err := userSvc.UserLogout(userId, platform)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).SuccessOk()
	return
}

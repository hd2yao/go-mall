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

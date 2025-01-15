package appservice

import (
	"context"
	"errors"

	"github.com/hd2yao/go-mall/api/reply"
	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

type UserAppSvc struct {
	ctx           context.Context
	userDomainSvc *domainservice.UserDomainSvc
}

func NewUserAppSvc(ctx context.Context) *UserAppSvc {
	return &UserAppSvc{
		ctx:           ctx,
		userDomainSvc: domainservice.NewUserDomainSvc(ctx),
	}
}

func (us *UserAppSvc) GenToken() (*reply.TokenReply, error) {
	token, err := us.userDomainSvc.GetAuthToken(12345678, "h5", "")
	if err != nil {
		return nil, err
	}
	logger.New(us.ctx).Info("generate token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	err = util.CopyProperties(tokenReply, token)
	if err != nil {
		err = errcode.Wrap("请求转换成 TokenReply 失败", err)
		return nil, err
	}
	return tokenReply, nil
}

func (us *UserAppSvc) TokenRefresh(refreshToken string) (*reply.TokenReply, error) {
	token, err := us.userDomainSvc.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	logger.New(us.ctx).Info("refresh token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	err = util.CopyProperties(tokenReply, token)
	if err != nil {
		err = errcode.Wrap("请求转换成 TokenReply 失败", err)
		return nil, err
	}
	return tokenReply, err
}

func (us *UserAppSvc) UserRegister(userRegisterReq *request.UserRegister) error {
	userInfo := new(do.UserBaseInfo)
	err := util.CopyProperties(userInfo, userRegisterReq)
	if err != nil {
		err = errcode.Wrap("请求转换成 UserBaseInfo 失败", err)
		return err
	}

	// 调用 domain service 注册用户
	_, err = us.userDomainSvc.RegisterUser(userInfo, userRegisterReq.Password)
	if errors.Is(err, errcode.ErrUserNameOccupied) {
		// 重名导致的注册不成功不需要额外处理
		return err
	}
	if err != nil && !errors.Is(err, errcode.ErrUserNameOccupied) {
		// TODO: 发通告告知用户注册失败 | | 记录日志、监控告警、提示有用户注册失败发生
		return err
	}

	// TODO: 写注册成功后的外围辅助逻辑，比如注册成功后给用户发确认邮件|短信

	// TODO: 如果产品逻辑是注册后帮用户登录，那这里再删掉登录的逻辑

	return nil
}

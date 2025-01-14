package domainservice

import (
	"context"
	"time"

	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/cache"
	"github.com/hd2yao/go-mall/logic/do"
)

type UserDomainSvc struct {
	ctx context.Context
}

func NewUserDomainSvc(ctx context.Context) *UserDomainSvc {
	return &UserDomainSvc{ctx: ctx}
}

// GetUserBaseInfo 获取用户基本信息(因为还没开发注册登录功能，先 Mock 一个返回)
func (us *UserDomainSvc) GetUserBaseInfo(userId int64) *do.UserBaseInfo {
	return &do.UserBaseInfo{
		ID:        12345678,
		NickName:  "hd2yao",
		LoginName: "hd2yao",
		Verified:  1,
		Avatar:    "",
		Slogan:    "",
		IsBlocked: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetAuthToken 生成 AccessToken 和 RefreshToken
// 在缓存中会存储最新的 Token 以及与 Platform 对应的 UserSession 同时会删除缓存中旧的 Token
func (us *UserDomainSvc) GetAuthToken(userId int64, platform string, sessionId string) (*do.TokenInfo, error) {
	// 获取用户基本信息
	user := us.GetUserBaseInfo(userId)
	// 处理异常情况：用户不存在、被删除、被禁用
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		err := errcode.ErrUserInvalid
		return nil, err
	}

	// 设置 userSession 缓存
	userSession := new(do.SessionInfo)
	userSession.UserId = userId
	userSession.Platform = platform
	if sessionId == "" {
		// 为空说明是用户的登录行为，重新生成 sessionId
		sessionId = util.GenSessionId(userId)
	}
	userSession.SessionId = sessionId
	accessToken, refreshToken, err := util.GenUserAuthToken(userId)
	if err != nil {
		err = errcode.Wrap("Token 生成失败", err)
		return nil, err
	}
	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken

	// 向缓存中设置 AccessToken 和 RefreshToken 的缓存
	err = cache.SetUserToken(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("设置 Token 缓存时发生错误", err)
		return nil, err
	}
	err = cache.DelOldSessionTokens(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("删除旧Token时发生错误", err)
		return nil, err
	}
	err = cache.SetUserSession(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("设置Session缓存时发生错误", err)
		return nil, err
	}

	srvCreateTime := time.Now()
	tokenInfo := &do.TokenInfo{
		AccessToken:   userSession.AccessToken,
		RefreshToken:  userSession.RefreshToken,
		Duration:      int64(enum.AccessTokenDuration.Seconds()),
		SrvCreateTime: srvCreateTime,
	}
	return tokenInfo, nil
}

package domainservice

import (
	"context"
	"time"

	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
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

func (us *UserDomainSvc) RefreshToken(refreshToken string) (*do.TokenInfo, error) {
	log := logger.New(us.ctx)
	// 获取锁以防止并发刷新
	ok, err := cache.LockTokenRefresh(us.ctx, refreshToken)
	defer cache.UnLockTokenRefresh(us.ctx, refreshToken)
	if err != nil {
		err = errcode.Wrap("刷新 Token 时设置 redis 锁发生错误", err)
		return nil, err
	}
	if !ok {
		err = errcode.ErrTooManyRequests
		return nil, err
	}

	// 获取 refreshToken 对应的缓存
	tokenSession, err := cache.GetRefreshToken(us.ctx, refreshToken)
	if err != nil {
		log.Error("GetRefreshTokenCacheErr", "err", err)
		// 服务端发生错误一律提示客户端 Token 有问题
		// 生产环境可以做好监控日志中这个错误的监控
		err = errcode.ErrToken
		return nil, err
	}

	// refreshToken 没有对应的缓存
	if tokenSession == nil || tokenSession.UserId == 0 {
		err = errcode.ErrToken
		return nil, err
	}

	// 获取用户在指定平台中的 Session 信息
	userSession, err := cache.GetUserPlatformSession(us.ctx, tokenSession.UserId, tokenSession.Platform)
	if err != nil {
		log.Error("GetUserPlatformSessionErr", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	// 请求刷新的 refreshToken 和缓存中的 refreshToken 不一致，证明这个 refreshToken 已经过时
	// RefreshToken 被窃取或者前端页面刷 Token 不是串行的互斥操作都有可能造成这个情况
	if userSession.RefreshToken != refreshToken {
		// 记一条警告日志
		log.Warn("ExpiredRefreshToken", "requestToken", refreshToken, "newToken", userSession.RefreshToken, "userId", userSession.UserId)
		// 错误返回 Token 不正确，或者更精细化的错误提示：已在xxx登录，如不是您本人操作请xxx
		err = errcode.ErrToken
		return nil, err
	}

	// 重新生成 Token 因为不是用户主动登录所以 sessionId 与之前保持一致
	tokenInfo, err := us.GetAuthToken(tokenSession.UserId, tokenSession.Platform, tokenSession.SessionId)
	if err != nil {
		err = errcode.Wrap("GenAuthTokenErr", err)
		return nil, err
	}
	return tokenInfo, nil
}

// VerifyAccessToken 验证 Token 是否有效
func (us *UserDomainSvc) VerifyAccessToken(accessToken string) (*do.TokenVerify, error) {
	tokenInfo, err := cache.GetAccessToken(us.ctx, accessToken)
	if err != nil {
		logger.New(us.ctx).Error("GetAccessTokenCacheErr", "err", err)
		return nil, err
	}
	tokenVerify := new(do.TokenVerify)
	if tokenInfo != nil && tokenInfo.UserId != 0 {
		tokenVerify.UserId = tokenInfo.UserId
		tokenVerify.SessionId = tokenInfo.SessionId
		tokenVerify.Approved = true
	} else {
		tokenVerify.Approved = false
	}
	return tokenVerify, nil
}

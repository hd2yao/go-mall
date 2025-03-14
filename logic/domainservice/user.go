package domainservice

import (
	"context"
	"time"

	"github.com/hd2yao/go-mall/api/request"
	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/cache"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/logic/do"
)

type UserDomainSvc struct {
	ctx     context.Context
	UserDao *dao.UserDao
}

func NewUserDomainSvc(ctx context.Context) *UserDomainSvc {
	return &UserDomainSvc{
		ctx:     ctx,
		UserDao: dao.NewUserDao(ctx),
	}
}

// GetUserBaseInfo 获取用户基本信息
func (us *UserDomainSvc) GetUserBaseInfo(userId int64) *do.UserBaseInfo {
	log := logger.New(us.ctx)
	user, err := us.UserDao.FindUserById(userId)
	if err != nil {
		log.Error("GetUserBaseInfoError", "err", err)
		return nil
	}
	userBaseInfo := new(do.UserBaseInfo)
	err = util.CopyProperties(userBaseInfo, user)
	if err != nil {
		log.Error(errcode.ErrCoverData.Msg(), "err", err)
		return nil
	}
	return userBaseInfo
}

// UpdateUserBaseInfo 更新用户基本信息
func (us *UserDomainSvc) UpdateUserBaseInfo(request *request.UserInfoUpdate, userId int64) error {
	user, err := us.UserDao.FindUserById(userId)
	if err != nil {
		return err
	}
	user.Nickname = request.Nickname
	user.Avatar = request.Avatar
	user.Slogan = request.Slogan
	err = us.UserDao.UpdateUser(user)
	return err
}

// GetAuthToken 生成 AccessToken 和 RefreshToken
// 在缓存中会存储最新的 Token 以及与 Platform 对应的 UserSession 同时会删除缓存中旧的 Token
func (us *UserDomainSvc) GetAuthToken(userId int64, platform string, sessionId string) (*do.TokenInfo, error) {
	// 获取用户基本信息
	user := us.GetUserBaseInfo(userId)
	// 确认是否为有效用户
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
	// 使用 refreshToken 更新时，sessionId 要保持一致
	userSession.SessionId = sessionId
	// 重新生成 AccessToken 和 RefreshToken
	accessToken, refreshToken, err := util.GenUserAuthToken(userId)
	if err != nil {
		err = errcode.Wrap("Token 生成失败", err)
		return nil, err
	}
	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken

	// 向缓存中设置新的 AccessToken 和 RefreshToken 的缓存
	err = cache.SetUserToken(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("设置 Token 缓存时发生错误", err)
		return nil, err
	}
	// 删除已有的 Session 缓存，包括 AccessToken 和 RefreshToken
	err = cache.DelOldSessionTokens(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("删除旧Token时发生错误", err)
		return nil, err
	}
	// 更新当前用户的 Session 缓存
	err = cache.SetUserSession(us.ctx, userSession)
	if err != nil {
		err = errcode.Wrap("设置Session缓存时发生错误", err)
		return nil, err
	}

	// 返回 Token 信息
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
		tokenVerify.Platform = tokenInfo.Platform
		tokenVerify.Approved = true
	} else {
		tokenVerify.Approved = false
	}
	return tokenVerify, nil
}

func (us *UserDomainSvc) RegisterUser(userInfo *do.UserBaseInfo, plainPassword string) (*do.UserBaseInfo, error) {
	// 确定登录名可用
	existedUser, err := us.UserDao.FindUserByLoginName(userInfo.LoginName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcRegisterUserError", err)
	}
	// 用户名已存在
	if existedUser.LoginName != "" {
		return nil, errcode.ErrUserNameOccupied
	}
	// 密码加密
	passwordHash, err := util.BcryptPassword(plainPassword)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}
	// 创建用户
	userModel, err := us.UserDao.CreateUser(userInfo, passwordHash)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}

	err = util.CopyProperties(userInfo, userModel)
	if err != nil {
		return nil, errcode.ErrCoverData
	}
	return userInfo, nil
}

func (us *UserDomainSvc) LoginUser(loginName string, plainPassword, platform string) (*do.TokenInfo, error) {
	existedUser, err := us.UserDao.FindUserByLoginName(loginName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcRegisterUserError", err)
	}
	if existedUser.ID == 0 {
		return nil, errcode.ErrUserNotRight
	}

	// 验证密码
	if !util.BcryptCompare(existedUser.Password, plainPassword) {
		return nil, errcode.ErrUserNotRight
	}

	// 生成 Token 和 Session
	tokenInfo, err := us.GetAuthToken(existedUser.ID, platform, "")
	return tokenInfo, err
}

func (us *UserDomainSvc) LogoutUser(userId int64, platform string) error {
	log := logger.New(us.ctx)
	userSession, err := cache.GetUserPlatformSession(us.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}

	// 删除用户的 AccessToken 和 RefreshToken
	err = cache.DelAccessToken(us.ctx, userSession.AccessToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	err = cache.DelRefreshToken(us.ctx, userSession.RefreshToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}

	// 删除用户在对应平台的 Session
	err = cache.DelUserSessionOnPlatform(us.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}

	return nil
}

// ApplyForPasswordReset 申请重置密码
func (us *UserDomainSvc) ApplyForPasswordReset(loginName string) (passwordResetToken, code string, err error) {
	user, err := us.UserDao.FindUserByLoginName(loginName)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordReset", err)
		return
	}
	if user.ID == 0 {
		err = errcode.ErrUserNotRight
		return
	}

	// 生成重置密码 Token 和验证码 code
	token, err := util.GenPasswordResetToken(user.ID)
	code = util.RandNumStr(6)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordReset", err)
		return
	}

	// 将 Token 和 code 存入缓存
	err = cache.SetPasswordResetToken(us.ctx, user.ID, token, code)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordReset", err)
		return
	}

	// TODO: 发送验证码 code 到用户邮箱或手机

	// 发送成功后返回 Token
	passwordResetToken = token
	return
}

func (us *UserDomainSvc) ResetPassword(resetToken, resetCode, newPlainPassword string) error {
	log := logger.New(us.ctx)
	userID, code, err := cache.GetPasswordResetToken(us.ctx, resetToken)
	if err != nil {
		log.Error("ResetPasswordError", "err", err)
		err = errcode.Wrap("ResetPasswordError", err)
		return err
	}

	// 验证 Token 和 code 是否匹配
	if userID == 0 || code != resetCode {
		return errcode.ErrParams
	}

	user, err := us.UserDao.FindUserById(userID)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}
	// 找不到用户或者用户为封禁状态
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		return errcode.ErrUserInvalid
	}

	newPass, err := util.BcryptPassword(newPlainPassword)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}
	// 更新用户密码
	user.Password = newPass
	err = us.UserDao.UpdateUser(user)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}

	// 删除用户所有存在的 Session
	err = cache.DelUserSession(us.ctx, user.ID)
	if err != nil {
		log.Error("ResetPasswordError", "err", err)
	}

	// 删除用户重置密码的 Token
	err = cache.DelPasswordResetToken(us.ctx, resetToken)
	if err != nil {
		log.Error("ResetPasswordError", "err", err)
	}
	return nil
}

// AddUserAddress 新增用户收货地址
func (us *UserDomainSvc) AddUserAddress(addressInfo *do.UserAddressInfo) (*do.UserAddressInfo, error) {
	addressModel, err := us.UserDao.CreateUserAddress(addressInfo)
	if err != nil {
		err = errcode.Wrap("AddUserAddress", err)
		return nil, err
	}
	err = util.CopyProperties(addressInfo, addressModel)
	if err != nil {
		return nil, errcode.ErrCoverData
	}
	return addressInfo, nil
}

// GetUserAddresses 查询用户收货信息列表
func (us *UserDomainSvc) GetUserAddresses(userId int64) ([]*do.UserAddressInfo, error) {
	addresses, err := us.UserDao.FindUserAddresses(userId)
	if err != nil {
		err = errcode.Wrap("GetUserAddresses", err)
		return nil, err
	}

	userAddresses := make([]*do.UserAddressInfo, 0)
	if len(addresses) == 0 {
		return userAddresses, nil
	}
	err = util.CopyProperties(&userAddresses, addresses)
	if err != nil {
		return nil, errcode.ErrCoverData
	}
	return userAddresses, nil
}

// GetUserSingleAddress 获取单个地址信息
func (us *UserDomainSvc) GetUserSingleAddress(userId int64, addressId int64) (*do.UserAddressInfo, error) {
	address, err := us.UserDao.GetSingleAddress(addressId)
	if err != nil || address.UserId != userId {
		logger.New(us.ctx).Error("UserAddressNotMatchError", "err", err, "return data", address, "addressId", addressId, "userId", userId)
		return nil, errcode.ErrParams
	}

	userAddress := new(do.UserAddressInfo)
	err = util.CopyProperties(userAddress, address)
	if err != nil {
		return nil, errcode.ErrCoverData
	}
	return userAddress, nil
}

// ModifyUserAddress 更改用户的地址信息
func (us *UserDomainSvc) ModifyUserAddress(address *do.UserAddressInfo) error {
	addressModel, err := us.UserDao.GetSingleAddress(address.ID)
	if err != nil || address.UserId != addressModel.UserId {
		// 不匹配的情况打印一条日志，监控系统按日志里的关键词做一下监控，好发现问题
		logger.New(us.ctx).Error("UserAddressNotMatchError", "err", err, "return data", addressModel, "request data", address)
		return errcode.ErrParams
	}
	err = us.UserDao.UpdateUserAddress(address)
	if err != nil {
		err = errcode.Wrap("UpdateUserAddressError", err)
	}
	return err
}

func (us *UserDomainSvc) DeleteOneUserAddress(userId, addressId int64) error {
	address, err := us.UserDao.GetSingleAddress(addressId)
	if err != nil || address.UserId != userId {
		logger.New(us.ctx).Error("UserAddressNotMatchError", "err", err, "return data", address, "addressId", addressId, "userId", userId)
		return errcode.ErrParams
	}
	err = us.UserDao.DeleteOneAddress(address)
	return err
}

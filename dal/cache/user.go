package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/hd2yao/go-mall/common/enum"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/logic/do"
)

// SetUserToken 设置用户 AccessToken 和 RefreshToken 缓存
func SetUserToken(ctx context.Context, session *do.SessionInfo) error {
	log := logger.New(ctx)
	err := setAccessToken(ctx, session)
	if err != nil {
		log.Error("redis error", "err", err)
		return err
	}
	err = setRefreshToken(ctx, session)
	if err != nil {
		log.Error("redis error", "err", err)
		return err
	}
	return nil
}

// 设置 AccessToken 缓存，以 AccessToken 为 key
func setAccessToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, session.AccessToken)
	sessionDataBytes, _ := json.Marshal(session)
	res, err := Redis().Set(ctx, redisKey, sessionDataBytes, enum.AccessTokenDuration).Result()
	logger.New(ctx).Debug("redis debug", "res", res, "err", err)
	return err
}

// 设置 RefreshToken 缓存，以 RefreshToken 为 key
func setRefreshToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, session.RefreshToken)
	sessionDataBytes, _ := json.Marshal(session)
	return Redis().Set(ctx, redisKey, sessionDataBytes, enum.RefreshTokenDuration).Err()
}

// SetUserSession 设置用户 Session 缓存
// 以 UserId 为 key，使用 Hash 存储
// Hash 中以 Platform 为 Key，存储相应用户的 Session 信息，这样多个平台登录后的会话不会相互干扰
func SetUserSession(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, session.UserId)
	sessionDataBytes, _ := json.Marshal(session)
	err := Redis().HSet(ctx, redisKey, session.Platform, sessionDataBytes).Err()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return err
	}
	return nil
}

// DelOldSessionTokens 删除用户旧 Session 的 Token
func DelOldSessionTokens(ctx context.Context, session *do.SessionInfo) error {
	oldSession, err := GetUserPlatformSession(ctx, session.UserId, session.Platform)
	if err != nil {
		return err
	}
	// 没有旧的 Session
	if oldSession != nil {
		return nil
	}
	err = DelAccessToken(ctx, oldSession.AccessToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	err = DelayDelRefreshToken(ctx, oldSession.RefreshToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return nil
}

// GetUserPlatformSession 获取用户在指定平台中的 Session 信息
func GetUserPlatformSession(ctx context.Context, userId int64, platform string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, userId)
	result, err := Redis().HGet(ctx, redisKey, platform).Result()
	if err != nil && err != redis.Nil { // redis.Nil 表示键或字段不存在
		return nil, err
	}
	// key 不存在
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	session := new(do.SessionInfo)
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// DelAccessToken 删除 AccessToken 缓存
func DelAccessToken(ctx context.Context, accessToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	return Redis().Del(ctx, redisKey).Err()
}

// DelayDelRefreshToken 刷新 Token 时让旧的 RefreshToken 保留一段时间自己过期
func DelayDelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	return Redis().Expire(ctx, redisKey, enum.OldRefreshTokenHoldingDuration).Err()
}

// DelRefreshToken 直接删除 RefreshToken 缓存(修改密码、退出登录时使用)
func DelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	return Redis().Del(ctx, redisKey).Err()
}

func LockTokenRefresh(ctx context.Context, refreshToken string) (bool, error) {
	redisLockKey := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().SetNX(ctx, redisLockKey, "locked", 10*time.Second).Result()
}

func UnLockTokenRefresh(ctx context.Context, refreshToken string) error {
	redisLockKey := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().Del(ctx, redisLockKey).Err()
}

func GetRefreshToken(ctx context.Context, refreshToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)
	return session, nil
}

func GetAccessToken(ctx context.Context, accessToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)
	return session, nil
}

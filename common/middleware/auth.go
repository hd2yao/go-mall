package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/hd2yao/go-mall/common/app"
	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/logic/domainservice"
)

// 用户认证相关的中间件

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("go-mall-token")
		// 生成的 token 是 40 个字符
		if len(token) != 40 {
			app.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}

		// 验证 token 是否有效
		tokenVerify, err := domainservice.NewUserDomainSvc(c).VerifyAccessToken(token)
		if err != nil {
			// 验证 Token 时服务出错
			app.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}
		if !tokenVerify.Approved {
			// Token 未通过验证
			app.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}

		c.Set("user_id", tokenVerify.UserId)
		c.Set("session_id", tokenVerify.SessionId)
		c.Set("platform", tokenVerify.Platform)
		c.Next()
	}
}

package middleware

import (
    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/common/util"
)

// 存放项目运行需要的基础中间件

func StartTrace() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceId := c.Request.Header.Get("traceid")
        pSpanId := c.Request.Header.Get("spanid")
        spanId := util.GenerateSpanID(c.Request.RemoteAddr)

        // 如果 traceId 为空，证明是链路的发端，把它设置成此次的 spanId
        if traceId == "" {
            // trace 标识整个请求的链路，span 标识链路中的不同服务
            traceId = spanId
        }

        c.Set("traceId", traceId)
        c.Set("spanId", spanId)
        c.Set("pSpanId", pSpanId)
        c.Next()
    }
}

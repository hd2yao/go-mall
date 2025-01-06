package middleware

import (
    "bytes"
    "io"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/hd2yao/go-mall/common/logger"
    "github.com/hd2yao/go-mall/common/util"
)

// 存放项目运行需要的基础中间件

// StartTrace 代码追踪中间件
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

/**
 * 程序无法直接从 gin 中拿到 response 的内容，所以需要自定义一个 ResponseWriter
 * 解决方案来自于 StackOverflow https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
 */

// 自定义 ResponseWriter，实现拦截和记录响应数据的功能
type bodyLogWriter struct {
    gin.ResponseWriter               // 嵌套 gin.ResponseWriter
    body               *bytes.Buffer // 用于保存响应内容
}

// Write 重写 Write 方法，拦截数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)                  // 将响应数据保存到 body 中
    return w.ResponseWriter.Write(b) // 继续写入到真正的响应流
}

// WriteString 重写 WriteString 方法，拦截字符串写入
func (w bodyLogWriter) WriteString(s string) (int, error) {
    w.body.WriteString(s)                  // 将响应数据保存到 body 中
    return w.ResponseWriter.WriteString(s) // 继续写入到真正的响应流
}

// LogAccess 日志记录中间件
func LogAccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 捕获请求体，保存 body
        // 读取 c.Request.Body 并保存到内存中，通过 io.NopCloser 重置了 Body，以便后续可以再次读取
        // go 1.16 之前使用 ioutil.ReadAll() 和 ioutil.NopCloser(), go 1.16 之后废弃
        // 现在使用 io.ReadAll() 和 io.NopCloser() 替换
        // 对于较大的请求体可能会导致内存消耗过高，可以设置最大读取大小
        // TODO:优化，把 body 保存到 context 中，而不是每次都重新读取
        // 尝试从 context 中获取 body，如果不存在则保存到 context 中
        if _, exists := c.Get("reqBody"); !exists {
            // 请求体未保存，读取并存储到 context
            reqBody, _ := io.ReadAll(io.LimitReader(c.Request.Body, 10<<20)) // 限制请求体大小
            c.Set("reqBody", reqBody)                                        // 保存到 context
            c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))          // 重置 Body
        }

        // 替换 ResponseWriter
        blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
        c.Writer = blw

        // 记录开始日志
        start := time.Now()
        accessLog(c, "access_start", time.Since(start), nil)
        defer func() {
            // 记录结束日志
            accessLog(c, "access_end", time.Since(start), blw.body.String())
        }()

        c.Next()
        return
    }
}

// accessLog 日志记录
func accessLog(c *gin.Context, accessType string, dur time.Duration, dataOut interface{}) {
    req := c.Request
    body, _ := c.Get("reqBody") // 从 context 获取请求体
    bodyStr := string(body.([]byte))
    query := req.URL.RawQuery
    path := req.URL.Path
    // TODO:实现 Token 认证后再把访问日志里也加上 token 记录
    // token := c.Request.Header.Get("token")

    logger.New(c).Info("AccessLog",
        "type", accessType,
        "ip", c.ClientIP(),
        // "token", token,
        "method", req.Method,
        "path", path,
        "query", query,
        "body", bodyStr,
        "output", dataOut,
        "time(ms)", int64(dur/time.Millisecond))
}

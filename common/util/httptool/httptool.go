package httptool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/logger"
	"github.com/hd2yao/go-mall/common/util"
)

var (
	_Client *http.Client
	once    sync.Once // 确保 HTTP Client 只被初始化一次，避免重复创建 http.Client，节省系统资源
)

func getHttpClient() *http.Client {
	if _Client != nil {
		// 因为 Unit test 里要把 Client 换掉，所以虽然用了 once.Do，但这里还是判断一下 _Client 有没有实例化
		return _Client
	}
	once.Do(func() {
		tr := &http.Transport{
			// Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,              // 最大空闲连接数
			IdleConnTimeout:       90 * time.Second, // 空闲连接超时 90 秒
			MaxIdleConnsPerHost:   50,
			MaxConnsPerHost:       50,
			ForceAttemptHTTP2:     true,             // 强制使用 HTTP/2
			TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时
			ExpectContinueTimeout: 1 * time.Second,  // `Expect: 100-Continue` 的超时时间
		}
		_Client = &http.Client{Transport: tr}
	})
	return _Client
}

// SetUTHttpClient 让单元测试能把 http Clint 覆盖成具有 Mock 拦截设置的 HttpClient
func SetUTHttpClient(client *http.Client) {
	_Client = client
}

func Request(method string, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	start := time.Now()
	reqOpts := defaultRequestOption() // 创建默认的 requestOption
	for _, opt := range options {     // 在 reqOpts 上应用所有选项
		if err = opt.apply(reqOpts); err != nil {
			return
		}
	}

	log := logger.New(reqOpts.ctx)
	defer func() {
		if err != nil {
			log.Error("HTTP_REQUEST_ERROR_LOG", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err)
		}
	}()

	// 创建请求对象
	req, err := http.NewRequest(method, url, bytes.NewReader(reqOpts.data))
	if err != nil {
		return
	}
	reqOpts.ctx, _ = context.WithTimeout(reqOpts.ctx, reqOpts.timeout) // 设置超时
	req = req.WithContext(reqOpts.ctx)
	defer req.Body.Close()

	// 在 Header 中添加追踪信息，把内部服务串起来
	traceId, spanId, _ := util.GetTraceInfoFromCtx(reqOpts.ctx)
	reqOpts.headers["traceId"] = traceId
	reqOpts.headers["spanId"] = spanId
	if len(reqOpts.headers) != 0 { // 设置请求头
		for k, v := range reqOpts.headers {
			req.Header.Add(k, v)
		}
	}

	// 发起请求
	client := getHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	defer func() {
		// 记录请求日志
		dur := time.Since(start).Milliseconds()
		if dur >= 3000 { // 超过 3秒 返回, 记一条 Warn 日志
			log.Warn("HTTP_REQUEST_SLOW_LOG", "method", method, "url", url, "body", string(reqOpts.data), "reply", string(respBody), "err", err, "time(ms)", dur)
		} else {
			log.Debug("HTTP_REQUEST_DEBUG_LOG", "method", method, "url", url, "body", string(reqOpts.data), "reply", string(respBody), "err", err, "time(ms)", dur)
		}
	}()

	httpStatusCode = resp.StatusCode
	if httpStatusCode != http.StatusOK {
		// 返回非 200 时Go的 http 库不回返回error, 这里处理成error 调用方好判断
		err = errcode.Wrap("request api error", errors.New(fmt.Sprintf("non 200 response, response code: %d", httpStatusCode)))
		return
	}

	respBody, _ = ioutil.ReadAll(resp.Body)
	return
}

// Get 发起 GET 请求
func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request("GET", url, options...)
}

// Post 发起 POST 请求
func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 默认自带 Header Content-Type: application/json 可通过传递 WithHeaders 增加或覆盖 Header 信息
	defaultHeader := map[string]string{"Content-Type": "application/json"}
	var newOptions []Option
	newOptions = append(newOptions, WithHeaders(defaultHeader), WithData(data), WithContext(ctx))
	newOptions = append(newOptions, options...)

	httpStatusCode, respBody, err = Request("POST", url, newOptions...)
	return
}

// 针对可选的 HTTP 请求配置项，模仿 gRPC 使用的 Options 设计模式实现
type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

// 创建默认的 requestOption
func defaultRequestOption() *requestOption {
	return &requestOption{
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: map[string]string{},
	}
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

// apply 方法实现
func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

// WithContext 配置选项：设置 Context
func WithContext(ctx context.Context) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.ctx = ctx
		return
	})
}

// WithTimeout 配置选项：设置 Timeout
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

// WithHeaders 配置选项：设置 Headers
func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return
	})
}

// WithData 配置选项：设置 Data
func WithData(data []byte) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}

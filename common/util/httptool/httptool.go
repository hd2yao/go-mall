package httptool

import (
    "context"
    "time"
)

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

// WithHeader 配置选项：设置 Headers
func WithHeader(headers map[string]string) Option {
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

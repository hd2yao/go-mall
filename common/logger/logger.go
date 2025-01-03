package logger

import (
    "context"
    "fmt"
    "path"
    "runtime"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type logger struct {
    ctx     context.Context
    traceId string
    spanId  string
    pSpanId string
    _logger *zap.Logger
}

func New(ctx context.Context) *logger {
    var traceId, spanId, pSpanId string
    if ctx.Value("traceId") != nil {
        traceId = ctx.Value("traceId").(string)
    }
    if ctx.Value("spanId") != nil {
        spanId = ctx.Value("spanId").(string)
    }
    if ctx.Value("pSpanId") != nil {
        pSpanId = ctx.Value("pSpanId").(string)
    }
    return &logger{
        ctx:     ctx,
        traceId: traceId,
        spanId:  spanId,
        pSpanId: pSpanId,
        _logger: _logger,
    }
}

func (l *logger) Debug(mag string, kv ...interface{}) {
    l.log(zap.DebugLevel, mag, kv...)
}

func (l *logger) Info(mag string, kv ...interface{}) {
    l.log(zap.InfoLevel, mag, kv...)
}

func (l *logger) Warn(mag string, kv ...interface{}) {
    l.log(zap.WarnLevel, mag, kv...)
}

func (l *logger) Error(mag string, kv ...interface{}) {
    l.log(zap.ErrorLevel, mag, kv...)
}

func (l *logger) log(lvl zapcore.Level, msg string, kv ...interface{}) {
    // 保证要打印的日志信息成对出现，默认补充一个 unknown 值
    if len(kv)%2 != 0 {
        kv = append(kv, "unknown")
    }
    // 日志信息中增加追踪参数
    kv = append(kv, "traceId", l.traceId, "spanId", l.spanId, "pSpanId", l.pSpanId)
    // 增加日志调用者信息，方便查日志时定位程序位置
    funcName, file, line := l.getLoggerCallerInfo()
    kv = append(kv, "funcName", funcName, "file", file, "line", line)

    // 将 kv 转换成 zap.Field 切片
    fields := make([]zap.Field, 0, len(kv)/2)
    for i := 0; i < len(kv); i += 2 {
        k := fmt.Sprintf("%v", kv[i])
        fields = append(fields, zap.Any(k, fmt.Sprintf("%v", kv[i+1])))
    }

    // 调用 Check 方法，判断这个日志级别是否能写入日志
    ce := l._logger.Check(lvl, msg)
    ce.Write(fields...)
}

// getLoggerCallerInfo 获取调用日志记录器的方法的函数名、文件名和行号
func (l *logger) getLoggerCallerInfo() (funcName, file string, line int) {
    // runtime.Caller(skip int)：skip 参数决定了要跳过的堆栈帧数量。
    // 0：当前函数（getLoggerCallerInfo）
    // 1：调用 getLoggerCallerInfo 的函数 function (l *logger) log(lvl zapcore.Level, msg string, kv ...interface{})
    // 2：再上一层的调用者 function (l *logger) Info(mag string, kv ...interface{})
    // 3：通常是用户直接调用日志记录器的方法，例如 logger.Info()
    pc, file, line, ok := runtime.Caller(3)
    if !ok {
        return "", "", 0
    }
    file = path.Base(file)
    funcName = runtime.FuncForPC(pc).Name()
    return funcName, file, line
}

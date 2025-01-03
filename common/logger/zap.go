package logger

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"

    "github.com/hd2yao/go-mall/common/enum"
    "github.com/hd2yao/go-mall/config"
)

var _logger *zap.Logger

// TODO: 暂时测试
func ZapLoggerTest(data interface{}) {
    _logger.Info("test for zap init",
        zap.Any("app", config.App),
        zap.Any("database", config.Database),
        zap.Any("data", "这是一个测试文本，目的是生成一个大约 200KB 大小的文本文件。"),
    )
}

func init() {
    // 为生产环境配置了一个标准的编码器配置
    encoderConfig := zap.NewProductionEncoderConfig()
    // 指定日志时间格式为 ISO8601
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    // 生成一个 JSON 编码器，用于将日志消息编码为 JSON 格式
    encoder := zapcore.NewJSONEncoder(encoderConfig)
    // 获取日志文件写入器，将日志写入文件
    fileWriteSyncer := getFileLogWriter()

    // 根据不同的环境配置日志输出
    var cores []zapcore.Core
    switch config.App.Env {
    case enum.ModeTest, enum.ModeProd:
        // 测试环境和生产环境的日志输出到文件中
        cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, zapcore.InfoLevel))
    case enum.ModeDev:
        // 开发环境同时向控制台和文件输出日志，Debug 级别的日志也会被输出
        cores = append(
            cores,
            zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
            zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
        )
    }

    // 使用 zapcore.NewTee 将多个日志输出器组合在一起，并使用这个 core 创建一个 logger
    core := zapcore.NewTee(cores...)
    _logger = zap.New(core)
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
    // 使用 lumberjack 实现 logger 的滚动
    lumberJackLogger := &lumberjack.Logger{
        Filename:   config.App.Log.FilePath,
        MaxSize:    config.App.Log.FileMaxSize,      // 文件最大 100M，当文件达到这个大小后 lumberjack 会自动进行切割，把原来的日志保存到备份文件中
        MaxBackups: config.App.Log.BackUpFileMaxAge, // 旧文件最多保留 60 天
        Compress:   false,                           // 是否压缩旧的日志文件
        LocalTime:  true,                            // 使用本地时间
    }
    return zapcore.AddSync(lumberJackLogger)
}

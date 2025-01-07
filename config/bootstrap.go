package config

import (
    "bytes"
    "embed"
    "os"

    "github.com/spf13/viper"
)

/**
 * 加载配置文件，把配置解析到配置对象中
 */

// 嵌入文件只能写在 embed 指令的 Go 文件的同级目录或者子目录中
//
//go:embed *.yaml
var configs embed.FS

func init() {
    env := os.Getenv("ENV")
    vp := viper.New()
    // 根据环境变量 ENV 决定要读取的应用启动配置
    configFileStream, err := configs.ReadFile("application." + env + ".yaml")
    if err != nil {
        panic(err)
    }

    vp.SetConfigType("yaml")
    err = vp.ReadConfig(bytes.NewBuffer(configFileStream))
    if err != nil {
        panic(err)
    }
    vp.UnmarshalKey("app", &App)

    vp.UnmarshalKey("database", &Database)
}

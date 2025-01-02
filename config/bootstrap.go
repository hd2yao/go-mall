package config

import (
    "os"
    "time"

    "github.com/spf13/viper"
)

/**
 * 加载配置文件，把配置解析到配置对象中
 */

const CONF_DIR = "config/"

func init() {
    env := os.Getenv("ENV")
    vp := viper.New()
    // 根据环境变量 ENV 决定要读取的应用启动配置
    configFilePath := CONF_DIR + "application.yaml"
    if env != "" {
        configFilePath = CONF_DIR + "application." + env + ".yaml"
    }
    vp.SetConfigFile(configFilePath)
    err := vp.ReadInConfig() // 查找并读取配置文件
    if err != nil {
        panic(err)
    }
    vp.UnmarshalKey("app", &App)

    vp.UnmarshalKey("database", &Database)
    Database.MaxLifeTime *= time.Second
}

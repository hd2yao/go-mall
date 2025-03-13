# go-mall

基于 Go 语言开发的电商系统

## 项目简介

go-mall 是一个基于 Go 语言开发的现代电商系统，采用了清晰的分层架构和最佳实践。

## 技术栈

- 框架：Gin
- 数据库：MySQL
- 缓存：Redis
- ORM：GORM
- 配置管理：Viper
- 日志：Zap
- 其他工具：copier, lo

## 项目结构

```Plain Text
.
├── api/            # API 层，包含路由、控制器、请求/响应定义
├── common/         # 公共代码
├── config/         # 配置文件
├── dal/           # 数据访问层
├── doc/           # 项目文档
├── event/         # 事件处理
├── library/       # 工具库
├── logic/         # 业务逻辑层
├── resources/     # 资源文件
├── storage/       # 存储相关
└── test/          # 测试代码
```

## 文档

详细文档请查看 `doc` 目录：

- [项目概述](doc/project-overview.md)
- [API 接口文档](doc/api/index.md)
- [数据库设计](doc/database-design.md)
- [系统架构](doc/architecture.md)

## 开发环境要求

- Go 1.20+
- MySQL 8.0+
- Redis 6.0+

## 快速开始

1.克隆项目

```bash
git clone https://github.com/hd2yao/go-mall.git
```

2.安装依赖

```bash
go mod download
```

3.配置环境

- 修改配置文件中的数据库和 Redis 连接信息

4.运行项目

```bash
go run main.go
```

## 许可证

本项目采用 MIT 许可证，详情请参见 [LICENSE](LICENSE) 文件

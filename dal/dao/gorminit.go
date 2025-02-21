package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/config"
)

var _DbMaster *gorm.DB
var _DbSlave *gorm.DB

// GORM V2 版本支持让一个自动按照执行的语句进行读写分离连接切换的功能 DBResolver
// https://gorm.io/zh_CN/docs/dbresolver.html

// 以下为手动读写切换的方法

// DB 返回只读实例
func DB() *gorm.DB {
	return _DbSlave
}

// DBMaster 返回主库实例
func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbMaster = initDB(config.Database.Master)
	_DbSlave = initDB(config.Database.Slave)
}

func getDialector(t, dsn string) gorm.Dialector {
	//switch t { 项目数据库需要加载多数据源时去掉注释
	//case "postgres":
	//	return postgres.Open(dsn)
	//default:
	//	return mysql.Open(dsn)
	//}
	return mysql.Open(dsn)
}

func initDB(option config.DbConnectOption) *gorm.DB {
	db, err := gorm.Open(
		getDialector(option.Type, option.DSN),
		&gorm.Config{ // 替换成本项目实现的 gormLogger
			Logger: NewGormLogger(),
		})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	// SetMaxOpenConns 设置数据库的最大打开连接数。
	sqlDB.SetMaxOpenConns(option.MaxOpenConn)
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(option.MaxIdleConn)
	// SetConnMaxLifetime 设置连接可重复使用的最长时间。
	sqlDB.SetConnMaxLifetime(option.MaxLifeTime)

	// 连接测试
	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}
	return db
}

// SetDBMasterConn 设置连接对象 -- 只用在单测中把DB连接改成sqlMock的DB连接
func SetDBMasterConn(conn *gorm.DB) {
	_DbMaster = conn
}

// SetDBSlaveConn 设置连接对象 -- 只用在单测中把DB连接改成sqlMock的DB连接
func SetDBSlaveConn(conn *gorm.DB) {
	_DbSlave = conn
}

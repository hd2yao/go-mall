package dao

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/dal/dao"
)

var (
	mock sqlmock.Sqlmock
	err  error
	db   *sql.DB
)

// TestMain 是在当前 package 下，最先运行的一个函数，常用于测试基础组件的初始化
func TestMain(m *testing.M) {
	// 这里创建一个 sqlmock 的数据库连接和 mock 对象，mock 对象管理 DB 预期要执行的 SQL 语句

	// sqlmock 默认使用 sqlmock.QueryMatcherRegex 作为默认的 SQL 匹配器
	// 该匹配器使用 mock.ExpectQuery 和 mock.ExpectExec 的参数作为正则表达式与真正执行的 SQL 语句进行匹配
	// 我们可以使用 regexp.QuoteMeta 把 SQL 转义成正则表达式 => mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`"))
	//
	// 如果想进行更严格的匹配, 可以让 sqlmock 使用 sqlmock.QueryMatcherEqual 作为匹配器匹配器
	// 该匹配器把 mock.ExpectQuery 和 mock.ExpectExec 的参数作为预期要执行的 SQL 语句跟真正要执行的 SQL 进行相等比较, 只有完全一样才会测试通过, 即使少个空格也不行
	// db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	// 把项目使用的 DB 替换成 sqlmock 的 DB 连接
	dbMasterConn, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
		DefaultStringSize:         0,
	}))
	dbSlaveConn, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
		DefaultStringSize:         0,
	}))
	dao.SetDBMasterConn(dbMasterConn)
	dao.SetDBSlaveConn(dbSlaveConn)

	// m.Run 是调用包下面各个 Test 函数的入口
	os.Exit(m.Run())
}

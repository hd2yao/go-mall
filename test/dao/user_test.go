package dao

import (
	"context"
	"database/sql"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/dal/dao"
	"github.com/hd2yao/go-mall/logic/do"
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

func TestUserDao_CreateUser(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	userInfo := &do.UserBaseInfo{
		NickName:  "Slang",
		LoginName: "slang@go-mall.com",
		Verified:  0,
		Avatar:    "",
		Slogan:    "happy!",
		IsBlocked: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	passwordHash, _ := util.BcryptPassword("123456")
	userIsDel := 0

	ud := dao.NewUserDao(context.TODO())
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(userInfo.NickName, userInfo.LoginName, passwordHash, userInfo.Verified, userInfo.Avatar,
			userInfo.Slogan, userIsDel, userInfo.IsBlocked, userInfo.CreatedAt, userInfo.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	userObj, err := ud.CreateUser(userInfo, passwordHash)
	assert.Nil(t, err)
	assert.Equal(t, userInfo.LoginName, userObj.LoginName)
}

package initialize

import (
	"fmt"
	"time"

	"crm/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 初始化MySQl数据库
// Mysql 根据配置信息初始化 MySQL 数据库连接。
// 该函数使用全局配置对象 global.Config.Mysql 中的参数来构建数据库连接字符串，
// 并使用 gorm 库打开数据库连接。此外，它还配置了数据库连接池的参数，
// 如最大空闲连接数、最大打开连接数和连接的最大生命周期。
// 最后，它将建立的数据库连接赋值给全局变量 global.Db。
func Mysql() {
	// 获取全局配置中的 MySQL 配置信息
	m := global.Config.Mysql
	// 构建 MySQL 数据库的 DSN (Data Source Name)
	s := "%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local"
	var dsn = fmt.Sprintf(s, m.Username, m.Password, m.Host, m.Port, m.Dbname)
	// 使用 gorm 开启 MySQL 数据库连接，并配置单个表的命名策略为单数形式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	// 如果打开数据库连接失败，则打印错误信息并返回
	if err != nil {
		fmt.Printf("mysql error: %s", err)
		return
	}
	// 获取底层的原始数据库连接
	sqlDb, err := db.DB()
	// 如果获取原始数据库连接失败，则只打印错误信息（不直接返回，因为全局数据库连接可能已被初始化）
	if err != nil {
		fmt.Printf("mysql error: %s", err)
	}
	// 配置数据库连接池的参数
	// 设置最大空闲连接数
	sqlDb.SetMaxIdleConns(m.MaxIdleConns)
	// 设置最大打开连接数
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	// 设置连接的最大生命周期
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
	// 将数据库连接赋值给全局变量
	global.Db = db
}

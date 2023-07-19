package mysql

import (
	"InterviewPush/models"
	"InterviewPush/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

var interviewRecord = models.InterviewRecord{}

func Init(cfg *settings.MySQLConfig) (err error) {
	//user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	fmt.Printf("%#v", dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
		return
	}
	//设置打开数据库连接的最大数量
	sqlDB, err := db.DB()
	if err != nil {
		panic("db.DB() failed")
		return
	}
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	db.AutoMigrate(&interviewRecord) //表迁移

	return
}

////持久化到数据库
//func Casbin() *casbin.Enforcer {
//	//casbin
//	// 使用 MySQL 数据库初始化一个 Xorm 适配器
//	a, err := xormadapter.NewAdapter("mysql", "root:123456@tcp(47.98.212.252:3306)/fever", true)
//	fmt.Println(err)
//	e, err := casbin.NewEnforcer("conf/rbac_models.conf", a)
//	fmt.Println(err)
//	//从DB加载策略
//	e.LoadPolicy()
//
//	return e
//}
//
////添加权限
//func AdCasbin(cm models.CasbinModel) bool {
//	e := Casbin()
//	add, _ := e.AddPolicy(cm.RoleName, cm.Path, cm.Method)
//	return add
//}

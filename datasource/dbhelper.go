package datasource

import (
	"fmt"
	"log"
	"sync"

	// MySQL驱动,也会实际加载进来,但是没有在该文件直接使用到,需要 _匿名 隐藏引入
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"lottery/conf"
)

// 互斥锁
var dbLock sync.Mutex

var masterInstance *xorm.Engine
var slaveInstance *xorm.Engine

// TODO - NOTICE 单例模式 - 得到唯一的主库实例
//  在应用运行期间会不断调用数据库操作,不能每次调用都实例化1次
func InstanceDbMaster() *xorm.Engine {
	// 如果存在直接返回,这里不需要锁
	if masterInstance != nil {
		return masterInstance
	}

	// 不存在要创建,创建前锁定
	dbLock.Lock()
	// 创建完成解锁
	defer dbLock.Unlock()

	// 锁定后 return NewDbMaster() 直接创建可能也会出问题
	// 有1个在创建,后面2个排队;这个创建完了后面2个再进来,又实例化,也不行
	// 导致失败单例

	// 还要再判断是否创建
	if masterInstance != nil {
		return masterInstance
	}

	return NewDbMaster()
}

// 返回xorm的MySQL数据库操作引擎
func NewDbMaster() *xorm.Engine {
	sourcename := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)

	instance, err := xorm.NewEngine(conf.DriverName, sourcename)
	if err != nil {
		log.Fatal("dbhelper.InstanceDbMaster NewEngine error ", err)
		return nil
	}
	// xorm支持的调试特性
	// SQL执行时间
	// instance.ShowExecTime()
	// 执行的SQL语句 生产false不展示 开发true展示
	instance.ShowSQL(true)
	//instance.ShowSQL(false)
	masterInstance = instance
	return masterInstance
}

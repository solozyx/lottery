package datasource

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"lottery/conf"
)

// 互斥锁
var rdsLock sync.Mutex

var cacheInstance *RedisConn

// 封装成一个redis资源池
type RedisConn struct {
	pool *redis.Pool
	// 设置是否打印redis日志
	showDebug bool
}

// 对外只有一个命令，封装了一个redis的命令
func (rds *RedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	// 从Redis连接池获取1个连接
	conn := rds.pool.Get()
	// 使用完该连接,需要把连接放回连接池
	defer conn.Close()

	// debug功能
	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		// 读取错误信息
		e := conn.Err()
		if e != nil {
			log.Println("rdshelper Do", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		// 微秒 us = 纳秒 / 1000
		fmt.Printf("[redis] [info] [%dus]cmd=%s, err=%s, args=%v, reply=%s\n",
			(t2-t1)/1000, commandName, err, args, reply)
	}
	return reply, err
}

// 设置是否打印操作日志
func (rds *RedisConn) ShowDebug(b bool) {
	// 私有字段showDebug的访问要提供公开方法
	rds.showDebug = b
}

// 单例模式 - 得到唯一的redis缓存实例
func InstanceCache() *RedisConn {
	if cacheInstance != nil {
		return cacheInstance
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	// 要有二次判断
	if cacheInstance != nil {
		return cacheInstance
	}
	return NewCache()
}

// RedisConn 实例化
func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				fmt.Sprintf("%s:%d", conf.RdsCache.Host, conf.RdsCache.Port))
			if err != nil {
				log.Fatal("rdshelper.NewCache Dial error ", err)
				return nil, err
			}
			// 创建Redis连接池成功 返回创建的连接
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		// 最大连接数
		MaxIdle: 10000,
		// 最大活跃数
		MaxActive: 10000,
		// 超时时间
		IdleTimeout: 0,
		// 等待
		Wait: false,
		// 最大连接活跃时间 一直活跃
		MaxConnLifetime: 0,
	}

	// 构建 RedisConn 对象
	instance := &RedisConn{
		pool: &pool,
	}
	// 单例赋值
	cacheInstance = instance
	// 打印redis执行命令日志
	cacheInstance.ShowDebug(true)
	//cacheInstance.ShowDebug(false)
	return cacheInstance
}

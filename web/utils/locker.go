/**
 * Redis实现 用户分布式锁
 */
package utils

import (
	"fmt"

	"lottery/datasource"
)

// 加锁 避免一个用户重复连续点击抽奖按钮,个人并发多次抽奖重入
func LockLucky(uid int) bool {
	return lockLuckyServ(uid)
}

// 解锁 用户调用抽奖接口完成及时释放锁,避免死锁
func UnlockLucky(uid int) bool {
	return unlockLuckyServ(uid)
}

// Redis实现分布式锁
func getLuckyLockKey(uid int) string {
	return fmt.Sprintf("lucky_lock_%d", uid)
}

func lockLuckyServ(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	// TODO : important 用Redis实现分布式锁
	// SET key = 1
	// NX 是否存在,不存在才能把 key 设置进去, 存在则不能设置进去;否则是更新key
	// EX 过期时间 3秒 执行该Redis操作,3秒钟锁还没有释放,根据过期时间自动释放
	// 过期时间是为了避免死锁,程序在运行中出现异常没有调用到 unLock 操作
	// 保证锁在3秒内能够释放
	// 过期时间不能太短,调用秒杀接口,在某个逻辑卡住了,处理不完这个请求,导致异常
	rs, _ := cacheObj.Do("SET", key, 1, "EX", 3, "NX")
	if rs == "OK" {
		return true
	} else {
		return false
	}
}

func unlockLuckyServ(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("DEL", key)
	if rs == "OK" {
		return true
	} else {
		return false
	}
}

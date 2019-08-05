/**
 * 同一个IP抽奖，每天的操作限制，本地或者redis缓存
 */
package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"lottery/comm"
	"lottery/datasource"
)

const ipFrameSize = 2

// init 程序启动时调用1次 清理ip分段数据
func init() {
	// 本地开发测试的时候，每次启动归零
	resetGroupIpList()
}

// TODO : 计划任务,定点在每日 0 时执行
//  重置单机IP今天次数
func resetGroupIpList() {
	log.Println("ip_day_lucky.resetGroupIpList start")
	cacheObj := datasource.InstanceCache()
	// 对2个段的数据 循环清理
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("ip_day_lucky.resetGroupIpList stop")

	// IP当天的统计数,次日零点清空,设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIpList)
}

// 今天的IP抽奖次数递增，返回递增后的数值
func IncrIpLuckyNum(strIp string) int64 {
	// string -> int 方便IP做散列处理 ,把1份大的数据,分为几份小数据
	ip := comm.Ip4toInt(strIp)
	// 和ip相关的数据散列为 2段 存储
	i := ip % ipFrameSize
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1)
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY error=", err)
		// 返回int32最大值 判断时就超过了设定的最大默认值
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}

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

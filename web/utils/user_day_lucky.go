package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"lottery/comm"
	"lottery/conf"
	"lottery/datasource"
)

const userFrameSize = 2

func init() {
	resetGroupUserList()
}

func resetGroupUserList() {
	log.Println("user_day_lucky.resetGroupUserList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < userFrameSize; i++ {
		key := fmt.Sprintf(conf.RdsDayUserLuckyCacheKeyPrefix + "%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("user_day_lucky.resetGroupUserList stop")

	// 用户当天的抽奖次数,次日零点归零
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList)
}

func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	key := fmt.Sprintf(conf.RdsDayUserLuckyCacheKeyPrefix + "%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1)
	if err != nil {
		log.Println("user_day_lucky redis HINCRBY key=", key,
			", uid = ", uid, ", error = ", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)
		return num
	}
}

func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf(conf.RdsDayUserLuckyCacheKeyPrefix + "%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, uid, num)
	if err != nil {
		log.Println("user_day_lucky redis HSET key=", key,
			", uid = ", uid, ", error = ", err)
	}
}

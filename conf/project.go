package conf

import "time"

// 同一个IP每天最多抽奖次数
const IpLimitMax = 500

// 用户每天最多抽奖次数
const UserPrizeMax = 1

// 同一个IP每天最多抽奖次数
const IpPrizeMax = 30

const GtypeVirtual = 0   // 虚拟币
const GtypeCodeSame = 1  // 虚拟券，相同的码
const GtypeCodeDiff = 2  // 虚拟券，不同的码
const GtypeGiftSmall = 3 // 实物小奖
const GtypeGiftLarge = 4 // 实物大奖

// TODO - NOTICE Go语言时间格式规则
const SysTimeform = "2006-01-02 15:04:05"
const SysTimeformShort = "2006-01-02"

// 是否需要启动全局计划任务服务
var RunningCrontabService = false

// 中国时区
var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")

// TODO - WARN 部署时根据实际情况更改该项
// ObjSalesign 签名密钥
var SignSecret = []byte("0123456789abcdef")
// cookie中的加密验证密钥
var CookieSecret = "hellolottery"

// Redis Cache Key
const RdsGiftCacheKey = "lottery:allgift"
const RdsUserCacheKeyPrefix = "lottery:info_user_"
const RdsBlackipCacheKeyPrefix = "lottery:info_blackip_"
const RdsDayIpLuckyCacheKeyPrefix = "lottery:day_ips_"
const RdsDayUserLuckyCacheKeyPrefix = "lottery:day_users_"
const RdsGiftPoolCacheKey = "lottery:gift_pool"
const RdsCodeCacheKeyPrefix = "lottery:gift_code_"

// 定义1天24小时奖品分配100%权重
var PrizeDataRandomDayTime = [100]int{
	// 24 * 3 = 72   平均3%的机会
	// 100 - 72 = 28 剩余28%的机会
	// 7 * 4 = 28    剩下的分别给7个小时增加4%的机会
	// 每个数字表示小时
	// 每个数字表示分配 1个机会
	0, 	0, 	0,
	1, 	1, 	1,
	2, 	2, 	2,
	3, 	3,	3,
	4, 	4, 	4,
	5, 	5, 	5,
	6, 	6, 	6,
	7, 	7, 	7,
	8, 	8, 	8,
	9, 	9, 	9, 		9, 	9, 	9, 	9,
	10, 10, 10, 	10, 10, 10, 10,
	11, 11, 11,
	12, 12, 12,
	13, 13, 13,
	14, 14, 14,
	15, 15, 15, 	15, 15, 15, 15,
	16, 16, 16, 	16, 16, 16, 16,
	17, 17, 17, 	17, 17, 17, 17,
	18, 18, 18,
	19, 19, 19,
	20, 20, 20, 	20, 20, 20, 20,
	21, 21, 21, 	21, 21, 21, 21,
	22, 22, 22,
	23, 23, 23,
}

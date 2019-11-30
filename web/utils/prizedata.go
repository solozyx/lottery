package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"

	"lottery/comm"
	"lottery/conf"
	"lottery/datasource"
	"lottery/models"
	"lottery/services"
)

func init() {
	// 本地开发测试的时候,每次重新启动,奖品池自动归零
	resetServGiftPool()
}

// 重置奖品池
func resetServGiftPool() {
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("DEL", conf.RdsGiftPoolCacheKey)
	if err != nil {
		log.Println("prizedata.resetServGiftPool DEL error = ", err)
	}
}

// 发奖 指定的奖品是否还可以发出来奖品
func PrizeGift(id, leftNum int) bool {
	ok := false
	ok = prizeServGift(id)
	if ok {
		// Redis奖品池发奖成功,更新MySQL数据库减少奖的库存
		giftService := services.NewGiftService()
		rows, err := giftService.DecrLeftNum(id, leftNum)
		// 影响行数小于1 表示没有更新
		if rows < 1 || err != nil {
			log.Println("prizedata.PrizeGift giftService.DecrLeftNum error = ", err,
				", rows = ", rows)
			// MySQL数据更新失败,不能发奖
			return false
		}
	}
	return ok
}

// 从Redis奖品池获取id奖品数量
func GetGiftPoolNum(id int) int {
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HGET", conf.RdsGiftPoolCacheKey, id)
	if err != nil {
		log.Println("prizedata.GetGiftPoolNum error = ", err)
		return 0
	}
	num := comm.GetInt64(rs, 0)
	return int(num)
}

// 奖品池发奖
func prizeServGift(id int) bool {
	cacheObj := datasource.InstanceCache()
	// 递减1 每次发1个奖品
	rs, err := cacheObj.Do("HINCRBY", conf.RdsGiftPoolCacheKey, id, -1)
	if err != nil {
		log.Println("prizedata.prizeServGift HINCRBY error = ", err)
		return false
	}
	// 默认值 -1
	num := comm.GetInt64(rs, -1)
	if num >= 0 {
		return true
	} else {
		return false
	}
}

// TODO : NOTICE 优惠券发放MySQL版本,废弃
// func PrizeCodeDiff_MySQL(giftId int, codeService services.CodeService) string {
func PrizeCodeDiffDeprecated(giftId int, codeService services.CodeService) string {
	// 使用负数,避免和正数的uid冲突
	lockUid := 0 - giftId - 100*1000*1000
	LockLucky(lockUid)
	defer UnlockLucky(lockUid)

	// 找到1个可用优惠券编码
	// 上次发放的优惠券编码
	codeId := 0
	codeInfo := codeService.NextUsingCode(giftId,codeId)
	if codeInfo != nil && codeInfo.Id > 0{
		// 发放优惠券 状态更新
		codeInfo.SysStatus = 2
		codeInfo.SysUpdated = comm.NowUnix()
		// 状态2不是空值 能更新成功 无需明确指定数据库更新字段
		codeService.Update(codeInfo,nil)
	} else {
		log.Println("prizedata.PrizeCodeDiff gift_id =  ",giftId)
		return ""
	}
	return codeInfo.Code
}

func PrizeCodeDiff(giftId int, codeService services.CodeService) string {
	return prizeServCodeDiff(giftId, codeService)
}

// 导入优惠券编码到Redis缓存
func ImportCacheCodes(id int, code string) bool {
	key := fmt.Sprintf(conf.RdsCodeCacheKeyPrefix + "%d", id)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.ImportCacheCodes SADD error = ", err)
		return false
	} else {
		return true
	}
}

// 重新整理优惠券编码到缓存
func RecacheCodes(id int, codeService services.CodeService) (succNum, failNum int) {
	// 从MySQL数据库查询该奖品对应的所有优惠券
	list := codeService.Search(id)
	if list == nil || len(list) <= 0 {
		return 0, 0
	}

	key := fmt.Sprintf(conf.RdsCodeCacheKeyPrefix + "%d", id)
	cacheObj := datasource.InstanceCache()
	// TODO : NOTICE 缓存已经有key了 做临时key
	//  比如MySQL手动删除了10个优惠券编码,这10个编码存在于正式key中
	//  从数据库查询这10个编码就没有,在正式key执行 SADD 操作 导致这 10个编码不会删除
	//  所以用新的 tmpKey 最后重命名 覆盖掉 正式 key
	tmpKey := "tmp_" + key
	for _, data := range list {
		// 优惠券正常状态 才写入缓存 ,已经发放 和 被删除的忽略
		if data.SysStatus == 0 {
			code := data.Code
			_, err := cacheObj.Do("SADD", tmpKey, code)
			if err != nil {
				log.Println("prizedata.RecacheCodes SADD error = ", err)
				failNum++
			} else {
				succNum++
			}
		}
	}
	// tmpKey 重命名为 key, 就把 key之前的数据覆盖更新了
	_, err := cacheObj.Do("RENAME", tmpKey, key)
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	}
	return succNum,failNum
}

func GetCacheCodeNum(id int, codeService services.CodeService) (int, int) {
	mysqlNum := 0
	redisNum := 0

	// 数据库中优惠券编码数量
	list := codeService.Search(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				mysqlNum++
			}
		}
	}

	// redis缓存中优惠券编码数量
	key := fmt.Sprintf(conf.RdsCodeCacheKeyPrefix + "%d", id)
	cacheObj := datasource.InstanceCache()
	// 统计 SCARD
	rs, err := cacheObj.Do("SCARD", key)
	if err != nil {
		log.Println("prizedata.GetCacheCodeNum SCARD error = ", err)
	} else {
		redisNum = int(comm.GetInt64(rs, 0))
	}

	return mysqlNum,redisNum
}

// 优惠券发放
func prizeServCodeDiff(id int, codeService services.CodeService) string {
	key := fmt.Sprintf(conf.RdsCodeCacheKeyPrefix + "%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SPOP", key)
	if err != nil {
		log.Println("prizedata.prizeServCodeDiff SPOP error = ", err)
		return ""
	}
	code := comm.GetString(rs, "")
	if code == "" {
		log.Printf("prizedata.prizeServCodeDiff rs = %s", rs)
		return ""
	}

	// 优惠券发放成功,更新数据库中发放状态
	codeService.UpdateByCode(&models.LtCode{
		Code:       code,
		SysStatus:  2,
		SysUpdated: comm.NowUnix(),}, nil)
	return code
}

// 奖品的发奖计划
func ResetGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	if giftInfo == nil || giftInfo.Id < 1 {
		return
	}
	id := giftInfo.Id
	nowTime := comm.NowUnix()

	// 不能发奖,不需要设置发奖周期
	if giftInfo.SysStatus == 1 || 			// 奖品状态不对
		giftInfo.TimeBegin >= nowTime || 	// 开始时间不对
		giftInfo.TimeEnd <= nowTime || 		// 结束时间不对
		giftInfo.LeftNum <= 0 || 			// 剩余数不足
		giftInfo.PrizeNum <= 0 { 			// 总数不足
		if giftInfo.PrizeData != "" {
			// 清空旧的发奖计划数据
			clearGiftPrizeData(giftInfo, giftService)
		}
		return
	}

	// 发奖周期
	dayNum := giftInfo.PrizeTime
	// 没有设置奖品发奖周期
	if dayNum <= 0 {
		// 不需要有发奖计划 把剩余奖品全部填充到奖品池,不用根据发奖计划
		setGiftPool(id, giftInfo.LeftNum)
		return
	}

	// 重置发奖计划数据 清空奖品池
	// 之前奖品池有剩余也不管 之前期的发奖计划
	// 只关心本期发奖计划
	setGiftPool(id, 0)

	// 发奖计划运算
	// [每天发奖数量]
	prizeNum := giftInfo.PrizeNum
	// 如奖品100个 2天 每天发奖50个
	avgNum := prizeNum / dayNum
	// 每天可以分配到的奖品数量 map[第x天]该天发奖数量
	dayPrizeNum := make(map[int]int)
	// 平均分配
	if avgNum >= 1 && dayNum > 0 {
		for day := 0; day < dayNum; day++ {
			dayPrizeNum[day] = avgNum
		}
	}
	// 剩下的随机分配到任意1天 如5个奖品 发2天 平均每天2个 剩余1个随机
	prizeNum -= dayNum * avgNum
	for prizeNum > 0 {
		prizeNum--
		// 在 dayNum 天随机取1天
		day := comm.Random(dayNum)
		_, ok := dayPrizeNum[day]
		if !ok {
			dayPrizeNum[day] = 1
		} else {
			dayPrizeNum[day] += 1
		}
	}
	// [每天的每分钟发奖数量]
	// 每天的map 每小时的map 60分钟的数组 奖品数量
	// map[第x天]map[第x小时][60分钟]发奖数量
	prizeData := make(map[int]map[int][60]int)
	for day, num := range dayPrizeNum {
		// 该天的发奖计划
		dayPrizeData := getGiftPrizeDataOneDay(num)
		prizeData[day] = dayPrizeData
	}
	// 数据结构 map[int]map[int][60]int 存储到MySQL数据库需要做json处理
	// 把该数据结构 格式化为 更简单的结构
	// 将周期内每天 每小时 每分钟的数据 prizeData 格式化,再json序列化保存到数据库表
	datalist := formatGiftPrizeData(nowTime, dayNum, prizeData)
	str, err := json.Marshal(datalist)
	if err != nil {
		log.Println("prizedata.ResetGiftPrizeData json.Marshal error = ", err)
	} else {
		// 保存奖品发奖计划到MySQL数据库
		info := &models.LtGift{
			Id:         giftInfo.Id,
			// 剩余数量 重置为 总数量
			LeftNum:    giftInfo.PrizeNum,
			PrizeData:  string(str),
			// 更新发奖周期开始时间为当前时间
			PrizeBegin: nowTime,
			// 到发奖周期的结束时间,需要重新计算发奖计划数据
			PrizeEnd:   nowTime + dayNum*86400,
			SysUpdated: nowTime,
		}
		err := giftService.Update(info, nil)
		if err != nil {
			log.Println("prizedata.ResetGiftPrizeData giftService.Update",
				info, ", error = ", err)
		}
	}
}

// 清空旧的奖品发放计划
func clearGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	info := &models.LtGift{
		Id:        giftInfo.Id,
		// 清空MySQL存储的奖品发奖计划
		PrizeData: "",
	}
	// xorm 要更新数据表字段为空字符串 要明确指定字段名称
	err := giftService.Update(info, []string{"prize_data"})
	if err != nil {
		log.Println("prizedata.clearGiftPrizeData giftService.Update",info, ", error = ", err)
	}
	setGiftPool(giftInfo.Id, 0)
}

// 设置奖品池 奖品id 奖品剩余数量
func setGiftPool(id, num int) {
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", conf.RdsGiftPoolCacheKey, id, num)
	if err != nil {
		log.Println("prizedata.setServGiftPool error = ", err)
	}
}

// 计算1天24小时的发奖计划
// 把num个数量奖品放到1天24小时,每个小时会出现不均匀发奖
func getGiftPrizeDataOneDay(num int) map[int][60]int {
	// 每小时 60分钟 发奖数量
	rs := make(map[int][60]int)
	// 每小时 发奖数量
	hourData := [24]int{}

	// [每个小时 发奖数量运算]
	// 将奖品分布到24个小时内
	if num > 100 {
		hourNum := 0
		// 奖品数量多 直接按照百分比计算
		for _, h := range conf.PrizeDataRandomDayTime {
			hourData[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourData[h]
			n := num * d / 100
			hourData[h] = n
			hourNum += n
		}
		num -= hourNum
	}

	// 奖品数量少的时候,或剩下了一些没有分配,需要用到随机概率来计算
	for num > 0 {
		num--
		// 通过随机数确定奖品落在哪个小时
		hourIndex := comm.Random(100)
		h := conf.PrizeDataRandomDayTime[hourIndex]
		hourData[h]++
	}

	// [每个小时 的 每个分钟 发奖数量运算]
	// 将每个小时内的奖品数量分配到60分钟
	for h, hnum := range hourData {
		if hnum <= 0 {
			continue
		}
		minuteData := [60]int{}
		if hnum >= 60 {
			avgMinute := hnum / 60
			for i := 0; i < 60; i++ {
				minuteData[i] = avgMinute
			}
			hnum -= avgMinute * 60
		}
		// 剩下数量随机到各分钟内
		for hnum > 0 {
			hnum--
			m := comm.Random(60)
			minuteData[m]++
		}
		// [每个小时 的 每个分钟 发奖数量]
		rs[h] = minuteData
	}
	return rs
}

// 当前时间 nowTime
// 发奖周期天数 dayNum
// 发奖数据 prizeData
// map[int]map[int][60]int ==> map[day][hour][minute]num
// 返回 [][2]int ==> [][发奖时间,发奖数量]
func formatGiftPrizeData(nowTime int, dayNum int, prizeData map[int]map[int][60]int) [][2]int {
	rs := make([][2]int, 0)
	nowHour := time.Now().Hour()

	// 日期发奖计划
	for dn := 0; dn < dayNum; dn++ {
		dayData, ok := prizeData[dn]
		if !ok {
			// 这1天没有发奖计划
			continue
		}
		// 小时发奖计划
		dayTime := nowTime + dn*86400
		for hn := 0; hn < 24; hn++ {
			// hourData,ok := dayData[hn]
			// 如当前是 18点 这样会读取到 0时数据,错乱,取模
			hourData, ok := dayData[(hn+nowHour)%24]
			if !ok {
				continue
			}
			// 分钟发奖计划
			hourTime := dayTime + hn*3600
			for mn := 0; mn < 60; mn++ {
				num := hourData[mn]
				if num <= 0 {
					continue
				}
				// 找到特定一个分钟的发奖计划数据
				minuteTime := hourTime + mn*60
				rs = append(rs, [2]int{minuteTime, num})
			}
		}
	}
	return rs
}

// 把所有奖品填充奖品池
func DistributionGiftPool() int {
	totalNum := 0
	now := comm.NowUnix()
	giftService := services.NewGiftService()
	// 读数据库,不读缓存,后台程序读取频率低
	list := giftService.GetAll(false)
	if list != nil && len(list) > 0 {
		for _, gift := range list {
			// 是否正常状态
			if gift.SysStatus != 0 {
				continue
			}
			// 是否限量产品
			if gift.PrizeNum < 1 {
				// 不限量的奖品,不做处理
				continue
			}
			// 时间段是否正常
			if gift.TimeBegin > now || gift.TimeEnd < now {
				continue
			}
			// 计划数据的长度太短,不需要解析和执行
			// 发奖计划 [[时间1,数量1],[时间2,数量2]]
			// 时间本身就超过7位字符串 PrizeData 不正常 也不做处理
			if len(gift.PrizeData) <= 7 {
				continue
			}
			var cronData [][2]int
			err := json.Unmarshal([]byte(gift.PrizeData), &cronData)
			if err != nil {
				log.Println("prizedata.DistributionGiftPool json.Unmarshal error = ", err)
			} else {
				index := 0
				giftNum := 0
				for i, data := range cronData {
					ct := data[0]
					num := data[1]
					if ct <= now {
						giftNum += num
						index = i + 1
					} else {
						// TODO : NOTICE 把发奖计划当前时间和之前时间奖品放入奖品池
						//  以后的时间就以后考虑
						break
					}
				}
				// 奖品放入到奖品池
				if giftNum > 0 {
					incrGiftPool(gift.Id, giftNum)
					totalNum += giftNum
				}
				// 有计划数据被执行过,需要更新到数据库
				if index > 0 {
					if index >= len(cronData) {
						cronData = make([][2]int, 0)
					} else {
						// index 之前的数据就没了
						cronData = cronData[index:]
					}
					str, err := json.Marshal(cronData)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool json.Marshal(cronData)",
							cronData, "error = ", err)
					}
					columns := []string{"prize_data"}
					err = giftService.Update(&models.LtGift{
						Id:        gift.Id,
						PrizeData: string(str),
					}, columns)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool giftService.Update error = ", err)
					}
				}
			}
		}
		if totalNum > 0 {
			// 超过1个奖品经过处理了,奖品发奖任务更新到数据库了,写入缓存
			// 数据库Update会清空缓存,这里手动做1次查询,重新建立缓存
			giftService.GetAll(true)
		}
	}
	return totalNum
}

// 根据计划数据 往奖品池增加奖品数量
func incrGiftPool(id, num int) int {
	cacheObj := datasource.InstanceCache()
	rtNum, err := redis.Int64(cacheObj.Do("HINCRBY", conf.RdsGiftPoolCacheKey, id, num))
	if err != nil {
		log.Println("prizedata.incrGiftPool error = ", err)
		return 0
	}

	if int(rtNum) < num {
		// 二次补偿
		num2 := num - int(rtNum)
		rtNum, err = redis.Int64(cacheObj.Do("HINCRBY", conf.RdsGiftPoolCacheKey, id, num2))
		if err != nil {
			log.Println("prizedata.incrGiftPool 2 error = ", err)
			return 0
		}
		// 三次补偿?  实际可能性非常低了,在几微秒redis出现大规模需要补偿情况概率非常小
		// 如果严格要求,做递归调用,不及预期一致调用下去
	}
	return int(rtNum)
}

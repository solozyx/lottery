package cron

import (
	"log"
	"time"

	"lottery/comm"
	"lottery/services"
	"lottery/web/utils"
)

func ConfigueAppOneCron() {
	go resetAllGiftPrizeData()
	go distributionAllGiftPool()
}

func resetAllGiftPrizeData() {
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	nowTime := comm.NowUnix()
	for _, giftInfo := range list {
		// 奖品设置了发奖周期
		if giftInfo.PrizeTime > 0 &&
			(giftInfo.PrizeData == "" || giftInfo.PrizeEnd <= nowTime) {
			// 计划任务在后台运行 日志拍错
			log.Println("crontab start utils.ResetGiftPrizeData giftInfo = ", giftInfo)
			// 更新数据库
			utils.ResetGiftPrizeData(&giftInfo, giftService)
			// 更新数据库后,缓存失效,读1遍缓存数据,重新建立缓存
			giftService.GetAll(true)
			log.Println("crontab end utils.ResetGiftPrizeData giftInfo = ", giftInfo)
		}
	}

	// 每5分钟执行一次
	time.AfterFunc(5 * time.Minute, resetAllGiftPrizeData)
}

func distributionAllGiftPool() {
	log.Println("crontab start utils.DistributionGiftPool")
	num := utils.DistributionGiftPool()
	log.Println("crontab end utils.DistributionGiftPool, num = ", num)

	// 每分钟执行一次
	time.AfterFunc(time.Minute, distributionAllGiftPool)
}

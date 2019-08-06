package controllers

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
)

// 最核心抽奖接口 GET http://localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1 验证登录用户
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录，再来抽奖"
		return rs
	}
	// 2 用户抽奖分布式锁定
	// TODO : important 用户抽奖分布式锁,抽奖接口调用是并发的,用户连续快速点击
	//  和网络发包请求刷抽奖接口
	//  分布式锁作用在用户上,对于全局的性能没有影响
	//  只是锁定某1个用户请求,而不是锁定抽奖接口
	ok := utils.LockLucky(loginuser.Uid)
	if ok {
		// 加锁成功必须解锁,防止死锁
		defer utils.UnlockLucky(loginuser.Uid)
	} else {
		rs["code"] = 102
		rs["msg"] = "正在抽奖 ... :) "
		return rs
	}

	// 3 验证用户今日参与抽奖次数
	userDayLuckyNum := utils.IncrUserLuckyNum(loginuser.Uid)
	// Redis缓存验证用户今日抽奖次数是否超限
	// TODO : NOTICE 如果redis重启,所有的用户执行递增加1操作, 0 + 1 = 1 则进入数据库验证分支
	//  但是数据库存储的用户今日参与抽奖次数可能大于1好多了
	if userDayLuckyNum > conf.UserPrizeMax {
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完,明天再来吧 :( "
		return  rs
	} else {
		// MySQL数据库验证用户今日抽奖次数是否超限
		// ok = c.checkUserday(loginuser.Uid)
		// TODO : NOTICE 这里把缓存获取的值传进去
		ok = c.checkUserday(loginuser.Uid,userDayLuckyNum)
		if !ok {
			rs["code"] = 103
			rs["msg"] = "今日抽奖次数已用完,明天再来吧 :( "
			return  rs
		}
	}

	// 4 验证 IP 今日参与抽奖次数
	ip := comm.ClientIP(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "今日相同IP参与抽奖次数已用完,明天再来吧 :( "
		return  rs
	}

	// 5 验证 IP 黑名单
	limitBlack := false
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	var blackIpInfo *models.LtBlackip
	if !limitBlack {
		ok,blackIpInfo = c.checkBlackip(ip)
		if !ok {
			fmt.Println("IP黑名单 ",ip,limitBlack)
			limitBlack = true
		}
	}

	// 6 验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack {
		ok,userInfo = c.checkBlackUser(loginuser.Uid)
		if !ok {
			fmt.Println("用户黑名单 ",loginuser.Uid,limitBlack)
			limitBlack = true
		}
	}

	// 7 给用户分配抽奖编码
	// 可设置精确率 1/10000
	// 随机数 [0,9999]
	prizeCode := comm.Random(10000)

	// 8 根据抽奖编码匹配是否中奖
	prizeGift := c.prize(prizeCode, limitBlack)
	if prizeGift == nil ||
		prizeGift.PrizeNum < 0 ||
		(prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		rs["code"] = 205
		rs["msg"] = "很遗憾，没有中奖，请下次再试 :( "
		return  rs
	}

	// 9 对数量有限制的奖品发奖校验
	if prizeGift.PrizeNum > 0 {
		ok = utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum)
		// 没有奖品可以发
		if !ok {
			rs["code"] = 207
			rs["msg"] = "很遗憾，没有中奖，请下次再试 :( "
			return  rs
		}
	}

	// 10 不同编码优惠券发放
	// 优惠券需要发放1个唯一编码,编码池没有可用编码,发奖失败
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode)
		if code == "" {
			rs["code"] = 208
			rs["msg"] = "很遗憾，没有中奖，请下次再试 :( "
			return  rs
		}
		prizeGift.Gdata = code
	}

	// 11 记录中奖信息保存
	result := models.LtResult{
		// Id自增长 无需赋值
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginuser.Uid,
		Username:   loginuser.Username,
		// 用户随机抽奖编码
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(&result)
	if err != nil {
		log.Println("index_lucky.GetLucky ServiceResult.Create ",
			result, ", error = ", err)
		rs["code"] = 209
		rs["msg"] = "很遗憾，没有中奖，请下次再试 :( "
		return  rs
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果获得了实物大奖,需要将用户 IP 设置成黑名单一段时间
		c.prizeLarge(ip, loginuser, userInfo, blackIpInfo)
	}

	// 12 返回抽奖结果
	rs["gift"] = prizeGift
	return rs

	//api := &LuckyApi{}
	//code, msg, gift := api.luckyDo(loginuser.Uid, loginuser.Username, ip)
	//rs["code"] = code
}

package controllers

import (
	"lottery/comm"
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
	// 4 验证 IP 今日参与抽奖次数
	// 5 验证 IP 黑名单
	// 6 验证用户黑名单
	// 7 给用户分配抽奖编码
	// 8 根据抽奖编码匹配是否中奖
	// 9 对数量有限制的奖品发奖校验
	// 10 不同编码优惠券发放
	// 11 记录中奖信息
	// 12 返回抽奖结果
	ip := comm.ClientIP(c.Ctx.Request())
	api := &LuckyApi{}
	code, msg, gift := api.luckyDo(loginuser.Uid, loginuser.Username, ip)
	rs["code"] = code
	rs["msg"] = msg
	rs["gift"] = gift
	return rs
}

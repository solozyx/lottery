package controllers

import (
	"lottery/conf"
	"lottery/models"
)

func (c *IndexController) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	// 获取可参与抽奖奖品列表
	giftList := c.ServiceGift.GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode && gift.PrizeCodeB >= prizeCode {
			// 中奖编码区间满足条件,说明该抽奖编码可以中奖
			// 非用户黑名单 非IP黑名单 非实物奖
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}

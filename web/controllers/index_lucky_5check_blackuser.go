package controllers

import (
	"time"

	"lottery/models"
	"lottery/services"
)

func (c *IndexController) checkBlackUser(uid int) (bool, *models.LtUser) {
	info := services.NewUserService().Get(uid)
	if info != nil && info.Blacktime > int(time.Now().Unix()) {
		// 黑名单存在并且有效 false表示不能验证通过
		return false, info
	}
	return true, info
}
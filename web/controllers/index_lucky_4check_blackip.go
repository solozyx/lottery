package controllers

import (
	"time"

	"lottery/models"
)

// 验证当前用户的IP是否存在黑名单限制
// 返回true表示该ip不在ip黑名单验证通过
func (c *IndexController) checkBlackip(ip string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIp(ip)
	if info == nil || info.Ip == "" {
		// 在MySQL中没有ip信息 表示该ip不是黑名单 验证通过 返回true
		return true, nil
	}
	// 黑名单期限大于当前时间
	if info.Blacktime > int(time.Now().Unix()) {
		// IP黑名单存在，而且没有洗白
		return false, info
	}
	return true, info
}

package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
)

// package 私有方法
func (c *IndexController) checkUserday(uid int, userDayLuckyNum int64) bool {
	userdayInfo := c.ServiceUserday.GetUserToday(uid)
	if userdayInfo != nil && userdayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userdayInfo.Num >= conf.UserPrizeMax {
			// TODO : NOTICE 重新初始化Redis缓存用户今日抽奖次数
			if int(userDayLuckyNum) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		} else {
			userdayInfo.Num++
			// TODO : NOTICE 重新初始化Redis缓存用户今日抽奖次数
			if int(userDayLuckyNum) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			err103 := c.ServiceUserday.Update(userdayInfo, nil)
			if err103 != nil {
				log.Println("index_lucky_check_userday ServiceUserDay.Update err103=", err103)
			}
		}
	} else {
		// 创建今天的用户参与抽奖记录
		y, m, d := time.Now().Date()
		// yyyymmdd 年是4位 月日是1位或者2位 统一做成2位
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userdayInfo = &models.LtUserday{
			// Id : 自增 不用赋值
			Uid:        uid,
			Day:        day,
			// 用户今天第1次参与抽奖
			Num:        1,
			// 创建时间
			SysCreated: int(time.Now().Unix()),
			// 更新时间 SysUpdated:
		}
		err103 := c.ServiceUserday.Create(userdayInfo)
		if err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Create err103=", err103)
		}

		// TODO : NOTICE 重新初始化Redis缓存用户今日抽奖次数
		utils.InitUserLuckyNum(uid, 1)
	}
	return true
}

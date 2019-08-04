package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"lottery/models"
	"lottery/services"
)

// 后台中奖记录
type AdminResultController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

// 中奖记录列表
func (c *AdminResultController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	uid := c.Ctx.URLParamIntDefault("uid", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	// 分页,每页 100 个数据
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	var datalist []models.LtResult
	if giftId > 0 {
		// 公共 gift 搜索
		datalist = c.ServiceResult.SearchByGift(giftId, page, size)
	} else if uid > 0 {
		// 通过 user 搜索
		datalist = c.ServiceResult.SearchByUser(uid, page, size)
	} else {
		// 没有gift和user搜索条件 则展示全部数据
		datalist = c.ServiceResult.GetAll(page, size)
	}
	// page默认值为1 如果查询结果不足1页 则总数就是该值
	total := (page - 1) + len(datalist)
	// 如果查询结果多余 100个 则查数据库确定数据总数
	if len(datalist) >= size {
		if giftId > 0 {
			// 通过奖品查询总数
			total = int(c.ServiceResult.CountByGift(giftId))
		} else if uid > 0 {
			// 通过用户查询总数
			total = int(c.ServiceResult.CountByUser(uid))
		} else {
			total = int(c.ServiceResult.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/result.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "result",
			"GiftId":   giftId,
			"Uid":      uid,
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

// 中奖记录删除
func (c *AdminResultController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceResult.Delete(id)
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

// 中奖记录作弊
func (c *AdminResultController) GetCheat() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		// 作弊 将 sys_status 设置为 2
		c.ServiceResult.Update(&models.LtResult{Id:id, SysStatus:2}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminResultController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		// 恢复 sys_status 设置为 0
		c.ServiceResult.Update(&models.LtResult{Id:id, SysStatus:0}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

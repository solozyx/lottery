package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"

	"lottery/services"
)

// Controller
type AdminController struct {
	Ctx iris.Context
	ServiceUser services.UserService
	ServiceGift services.GiftService
	ServiceCode services.CodeService
	ServiceResult services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

// Action
func (c *AdminController) Get() mvc.Result {
	// 返回 mvc 的模板对象
	return mvc.View{
		// 后台默认首页
		Name: "admin/index.html",
		// 模板数据
		Data: iris.Map{
			// 页面标题
			"Title": "管理后台",
			// 当前频道 默认首页控值即可
			"Channel":"",
		},
		// 布局
		Layout: "admin/layout.html",
	}
}

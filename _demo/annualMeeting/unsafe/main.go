/**
 * 年会抽奖程序
 * 不是线程安全
 * 基础功能：
 * 1 /import 导入参与名单作为抽奖的用户
 * 2 /lucky 从名单中随机抽取用户
 * 测试方法：
 * curl http://localhost:8080/
 * curl --data "users=yifan,yifan2" http://localhost:8080/import
 * curl http://localhost:8080/lucky
 */

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// 用户列表
var userList []string

// 抽奖控制器
type lotteryController struct {
	Ctx iris.Context
}

// 启动iris应用
func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	userList = make([]string, 0)
	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

// GET http://localhost:8080/
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d\n", count)
}

// POST http://localhost:8080/import
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	// 逗号分隔
	users := strings.Split(strUsers, ",")
	count1 := len(userList)
	for _, u := range users {
		// 去掉前后空白字符
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			// 导入用户
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d，成功导入用户数: %d\n", count2, (count2 - count1))
}

// 抽奖
// GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	if count > 1 {
		// 随机数种子 纳秒时间戳
		seed := time.Now().UnixNano()
		// 生成指定长度在 count 范围内的随机数 [0,count]
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		// 抽取到一个用户
		user := userList[index]
		// 移除这个用户
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = userList[0:0]
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("已经没有参与用户，请先通过 /import 导入用户 \n")
	}
}

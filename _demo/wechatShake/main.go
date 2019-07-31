/**
 * 微信摇一摇
 *
 * 基础功能：
 * /lucky 只有一个抽奖的接口，奖品信息都是预先配置好的
 * 测试方法：
 * curl http://localhost:8080/
 * curl http://localhost:8080/lucky
 * 压力测试：（线程不安全的时候，总的中奖纪录会超过总的奖品数）
 * wrk -t10 -c10 -d5 http://localhost:8080/lucky
 */

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// 奖品类型，枚举 iota 从0开始自增
const (
	giftTypeCoin      = iota // 虚拟币
	giftTypeCoupon           // 优惠券，不相同的编码
	giftTypeCouponFix        // 优惠券，相同的编码
	giftTypeRealSmall        // 实物小奖
	giftTypeRealLarge        // 实物大奖
)

// 文件日志
var logger *log.Logger

// 最大中奖号码
const rateMax = 10000

// 奖品信息
type gift struct {
	id       int      // 奖品ID
	name     string   // 奖品名称
	pic      string   // 照片链接
	link     string   // 链接
	gtype    int      // 奖品类型
	data     string   // 奖品的数据（特定的配置信息，如：虚拟币面值，固定优惠券的编码）
	datalist []string // 奖品数据集合（特定的配置信息，如：不同的优惠券的编码）
	total    int      // 总数，0 不限量
	left     int      // 剩余数
	inuse    bool     // 是否参与抽奖活动
	rate     int      // 中奖概率，万分之N,0-9999
	rateMin  int      // 大于等于，中奖的最小号码,0
	rateMax  int      // 小于，中奖的最大号码,9999
}

// 奖品列表
var giftlist []*gift

// 对共享变量 giftlist 加互斥锁
var mu sync.Mutex

// 抽奖的控制器
type lotteryController struct {
	Ctx iris.Context
}

// 初始化日志信息
func initLog() {
	// "/var/log/lottery_demo.log"
	f, _ := os.Create("./lottery_demo.log")
	// 日志前缀 日期+毫秒数
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

// 初始化奖品列表信息（管理后台来维护）
func initGift() {
	giftlist = make([]*gift, 5)
	// 1 实物大奖
	// 测试 total:1000, left:1000, rate:10000, 表示100%中奖
	// 测试 total:10, left:10, rate:10000, 表示10/10000中奖
	g1 := gift{
		id:       1,
		name:     "手机大奖",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealLarge,
		data:     "",
		datalist: nil,
		total:    2,
		left:     2,
		inuse:    true,
		rate:     1, // 表示1/10000中奖
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[0] = &g1

	// 2 实物小奖
	g2 := gift{
		id:       2,
		name:     "充电器",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealSmall,
		data:     "",
		datalist: nil,
		total:    5,
		left:     5,
		inuse:    true,
		rate:     10, // 表示中奖概率 10/10000
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[1] = &g2

	// 3 虚拟券，相同的编码
	g3 := gift{
		id:       3,
		name:     "优惠券满200元减50元",
		pic:      "",
		link:     "",
		gtype:    giftTypeCouponFix,
		data:     "mall-coupon-2019",
		datalist: nil,
		total:    50,
		left:     50,
		rate:     500, // 中奖概率 500/10000
		inuse:    true,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[2] = &g3

	// 4 虚拟券，不相同的编码
	g4 := gift{
		id:    4,
		name:  "商城无门槛直降50元优惠券",
		pic:   "",
		link:  "",
		gtype: giftTypeCoupon,
		data:  "",
		// TODO - NOTICE 注意这里要设置 total=left=len(datalist) 防止越界崩溃
		datalist: []string{"c01", "c02", "c03", "c04", "c05"},
		total:    5,
		left:     5,
		inuse:    true,
		rate:     100,
		rateMin:  0,
		rateMax:  0,
	}
	giftlist[3] = &g4

	// 5 虚拟币
	g5 := gift{
		id:      5,
		name:    "社区10个金币",
		pic:     "",
		link:    "",
		gtype:   giftTypeCoin,
		data:    "10金币",
		total:   5,
		left:    5,
		inuse:   true,
		rate:    5000,
		rateMin: 0,
		rateMax: 0,
	}
	giftlist[4] = &g5

	// 整理奖品数据，把rateMin,rateMax根据rate进行编排
	rateStart := 0
	for _, data := range giftlist {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = data.rateMin + data.rate
		// TODO - NOTICE 中奖概率的1个小算法
		if data.rateMax >= rateMax {
			// 号码达到最大值，分配的范围重头再来
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
	}
	fmt.Printf("giftlist=%v\n", giftlist)
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	// 初始化日志信息
	initLog()
	// 初始化奖品信息
	initGift()
	return app
}

func main() {
	app := newApp()
	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

// 奖品数量信息列表 GET http://localhost:8080/
func (c *lotteryController) Get() string {
	count := 0
	total := 0
	for _, data := range giftlist {
		// total=0表示不限制抽奖数量
		if data.inuse && (data.total == 0 ||
			// 有库存
			(data.total > 0 && data.left > 0)) {
			count++
			total += data.left
		}
	}
	return fmt.Sprintf("当前有效奖品种类数量: %d，限量奖品总数量=%d\n", count, total)
}

// 抽奖接口 GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() map[string]interface{} {
	// 互斥锁方案解决抽奖非线程安全
	mu.Lock()
	defer mu.Unlock()

	// 每个用户分配1个参与抽奖的编码
	code := luckyCode()
	ok := false
	result := make(map[string]interface{})
	result["success"] = ok
	// 抽奖逻辑
	for _, data := range giftlist {
		// 该产品抽奖不参与抽奖 或 库存不足
		if !data.inuse || (data.total > 0 && data.left <= 0) {
			continue
		}
		if data.rateMin <= int(code) && data.rateMax > int(code) {
			// 中奖了，抽奖编码在奖品中奖编码范围内
			// 根据奖品类型 调用不同发奖方法
			sendData := ""
			switch data.gtype {
			case giftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}
			// 中奖后，成功得到奖品（发奖成功）生成中奖纪录
			if ok {
				saveLuckyData(code, data.id, data.name, data.link, sendData, data.left)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = sendData
				break
			}
		}
	}
	return result
}

// 抽奖编码 0-9999
func luckyCode() int32 {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return code
}

// 发奖阶段
// 发奖，虚拟币
func sendCoin(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		// 还有剩余
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，优惠券（不同值）
func sendCoupon(data *gift) (bool, string) {
	if data.left > 0 {
		// 还有剩余的奖品
		left := data.left - 1
		data.left = left
		// 奖品初始值 datalist: []string{"c01", "c02", "c03", "c04", "c05"}
		return true, data.datalist[left]
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，优惠券（固定值）
func sendCouponFix(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，小实物
func sendRealSmall(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 发奖，大实物
func sendRealLarge(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left--
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

// 记录用户的获奖记录
func saveLuckyData(code int32, id int, name, link, sendData string, left int) {
	logger.Printf("lucky中奖了, code抽奖码=%d, gift奖品id=%d, name奖品名称=%s, link奖品链接=%s, data=%s, left奖品剩余=%d ",
		code, id, name, link, sendData, left)
}

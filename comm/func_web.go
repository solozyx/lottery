package comm

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"lottery/conf"
	"lottery/models"
)

const cookieName = "lottery_loginuser"

// 得到客户端IP地址
func ClientIP(request *http.Request) string {
	// 用户 IP:PORT
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

// 跳转URL
func Redirect(writer http.ResponseWriter, url string) {
	// 头信息 Location跳转到url
	writer.Header().Add("Location", url)
	// 301是永久重定向
	// 302是临时重定向,这里用 302
	writer.WriteHeader(http.StatusFound)
}

// 从cookie中得到当前登录的用户
func GetLoginUser(request *http.Request) *models.ObjLoginuser {
	c, err := request.Cookie(cookieName)
	if err != nil {
		// cookie不存在
		return nil
	}
	// cookie保存在url格式中 需要解析出cookie
	params, err := url.ParseQuery(c.Value)
	if err != nil {
		// 解析cookie错误
		return nil
	}
	// url是字符串类型
	uid, err := strconv.Atoi(params.Get("uid"))
	if err != nil || uid < 1 {
		// uid错误
		return nil
	}
	// Cookie最长使用时长
	now, err := strconv.Atoi(params.Get("now"))
	if err != nil || NowUnix()-now > 86400*30 {
		// 种在客户端的cookie超过30天,认为cookie失效
		return nil
	}

	// IP修改了是不是要重新登录
	ip := params.Get("ip")
	if ip != ClientIP(request) {
		return nil
	}

	// 构建登录对象
	loginuser := &models.ObjLoginuser{}
	loginuser.Uid = uid
	loginuser.Username = params.Get("username")
	loginuser.Now = now
	loginuser.Ip = ClientIP(request)
	// 签名信息
	loginuser.Sign = params.Get("sign")
	//if err != nil {
	//	log.Println("fuc_web GetLoginUser Unmarshal ", err)
	//	return nil
	//}

	// 验证客户端签名
	sign := createLoginuserSign(loginuser)
	// 验证客户端签名失败
	if sign != loginuser.Sign {
		log.Printf("fuc_web GetLoginUser FAIL : user cookie sign = %s and correct sign = %s",
			sign, loginuser.Sign)
		return nil
	}
	// TODO : 更新Cookie字段 now sign
	return loginuser
}

// 将登录成功的用户信息种到用户客户端浏览器cookie中
func SetLoginuser(writer http.ResponseWriter, loginuser *models.ObjLoginuser) {
	if loginuser == nil || loginuser.Uid < 1 {
		// 清理cookie 退出登录
		c := &http.Cookie{
			// cookie名称
			Name: cookieName,
			// cookie值
			Value: "",
			// 根目录
			Path: "/",
			// 让cookie过期
			MaxAge: -1,
		}
		http.SetCookie(writer, c)
		return
	}

	// 生成签名
	if loginuser.Sign == "" {
		loginuser.Sign = createLoginuserSign(loginuser)
	}

	// 构造cookie
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginuser.Uid))
	params.Add("username", loginuser.Username)
	params.Add("now", strconv.Itoa(loginuser.Now))
	params.Add("ip", loginuser.Ip)
	// 在cookie中加入签名,才能保证cookie可信
	// 客户端怎么上传数据无法控制,只能根据签名识别伪造请求
	params.Add("sign", loginuser.Sign)
	c := &http.Cookie{
		Name:  cookieName,
		Value: params.Encode(),
		Path:  "/",
	}
	http.SetCookie(writer, c)
}

// 根据登录用户信息生成加密字符串
func createLoginuserSign(loginuser *models.ObjLoginuser) string {
	// secret 字符串拼接规则 外部用户不知道
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s&now=%d",
		loginuser.Uid, loginuser.Username, conf.CookieSecret, loginuser.Now)
	return CreateSign(str)
}

// 对字符串进行签名
func CreateSign(str string) string {
	str = string(conf.SignSecret) + str
	// md5 不可逆
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}

/**
 * redis缓存配置信息
 */
package conf

type RdsConfig struct {
	Host      string
	Port      int
	User      string
	Pwd       string
	IsRunning bool // 是否正常运行
}

// 系统中用到的所有redis缓存资源
var RdsCacheList = []RdsConfig{
	{
		Host:      "192.168.174.134",
		Port:      6379,
		User:      "",
		Pwd:       "",
		IsRunning: true,
	},
}

var RdsCache RdsConfig = RdsCacheList[0]

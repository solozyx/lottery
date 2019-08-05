package models

// 站点中与浏览器交互的用户模型
type ObjLoginuser struct {
	// 用户唯一id
	Uid      int
	Username string
	// 时间戳
	Now      int
	Ip       string
	// 签名,签名生成 验证 cookie识别 序列化保存
	Sign     string
}

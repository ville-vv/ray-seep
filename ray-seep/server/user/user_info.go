package user

// 用户账户信息
type UserAccount struct {
	UserID   int64  // 用户ID号
	UserName string // 用户姓名
	AppKey   string // 用户的 app key
	Secret   string // 用户秘钥
}

// 用户协议信息
type UserProtocol struct {
	UserID       int64  `json:"user_id"`       // 用户ID号
	ProtocolName string `json:"protocol_name"` //  协议名称
	ProtocolPort int    `json:"protocol_port"` // 协议端口号
}

type UserLoginInfo struct {
	UserID     int64  // 用户ID号
	UserName   string // 用户姓名
	ServerKind string // 用户服务类型
	AppKey     string // 用户的 app key
	Secret     string // 用户秘钥
	Token      string // 用户登录标识
}

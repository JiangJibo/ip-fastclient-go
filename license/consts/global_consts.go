package consts

var (
	//不限制qps的qps值为-1
	MgUnlimitedQps = "-1"

	//cho模块aes加密key
	MgLsnEchoKey = "rFQd]2Xwl{Zy&Ex%"

	// Echo模块分隔符
	MgLsnEchoSep = "^|^"

	//rsa公钥前面追加的随机字符串长度
	MgPrefixConfusionSize = 76

	//rsa公钥后面追加的随机字符串长度
	MgPostfixConfusionSize = 76

	// aes密码长度, 长度必须为16（AES-128）
	MgAesPasswordSize = 16

	//系统时间和license颁发时间最多相差的时间是2天
	MaxDeltaSeconds = 2 * 24 * 3600

	// 混淆的chunk个数
	MaxChaosParts = 11

	//魔法数字
	MagicNum = 12

	// 到期30天前开始提醒
	MaxNotifyBeforeSeconds = 30 * 24 * 3600
)

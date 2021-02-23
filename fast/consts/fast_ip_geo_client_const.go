package consts

var (
	//SDK版本
	VERSION = "1.0.0"

	//元数据的字节长度
	META_INFO_BYTE_LENGTH = 1024

	//UTF-8编码
	CONTENT_CHAR_SET = "utf-8"

	//geo地理位置信息某个字段为空时候的返回值
	NOTFOUND_GEO_ITEM_VALUE = ""

	//分隔符，不可打印字符
	GEO_RAW_SEP = "\u0000"

	//经度字段名，做水印用
	GEO_X = "longitude"

	//纬度字段名，做水印用
	GEO_Y = "latitude"

	//混淆字符串的长度
	MG_CONFUSED_SIZE = 2048

	//RC4 加密key的长度, key不能太长，不然用户必须要安装jce, 太麻烦 https://stackoverflow.com/questions/6481627/java-security-illegal-key-size-or-default-parameters)
	MG_KEY_SIZE = 16

	//前缀混淆长度
	MG_KEY_START_INDEX = 714

	//文件加密算法
	MG_DAT_ALGORITHM = "RC4"
)

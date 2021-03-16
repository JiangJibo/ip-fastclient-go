package consts

var (
	//SDK版本
	VERSION = "1.0.0"

	//元数据的字节长度
	MetaInfoByteLength = 1024

	//UTF-8编码
	ContentCharSet = "utf-8"

	//geo地理位置信息某个字段为空时候的返回值
	NotfoundGeoItemValue = ""

	//分隔符，不可打印字符
	GeoRawSep = "\u0000"

	//经度字段名，做水印用
	GeoX = "longitude"

	//纬度字段名，做水印用
	GeoY = "latitude"

	//混淆字符串的长度
	MgConfusedSize = 2048

	//RC4 加密key的长度, key不能太长，不然用户必须要安装jce, 太麻烦 https://stackoverflow.com/questions/6481627/java-security-illegal-key-size-or-default-parameters)
	MgKeySize = 16

	//前缀混淆长度
	MgKeyStartIndex = 714

	//文件加密算法
	MgDatAlgorithm = "RC4"
)

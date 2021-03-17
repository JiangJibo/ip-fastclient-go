package impl

import (
	"ip-fastclient-go/fast/consts"
	"ip-fastclient-go/fast/context"
	"ip-fastclient-go/fast/xprt"
	"ip-fastclient-go/license/client"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

var (
	// ipv6 a段索引区占用字节数
	Ipv6FirstSegmentSize = 256 * 256
)

// 保留ip段
type ReservedIpRange struct {
	start        []byte
	end          []byte
	contentIndex int
}

type ipRanges []*ReservedIpRange

//调用标准库的sort.Sort必须要先实现Len(),Less(),Swap() 三个方法.
func (s ipRanges) Len() int { return len(s) }

func (s ipRanges) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ipRanges) Less(i, j int) bool { return len(s[i].end) < len(s[j].end) }

type IPv6GeoClient struct {
	// ipBlockStartIndex[1][128.100] = 1000 : 表示以128.100开头第一个IP段的序号是1000
	ipBlockStartIndex [][]uint32

	// ipBlockStartIndex[128.100] = 1200 : 表示以128.100开始最后一个IP段的序号是1200
	ipBlockEndIndex [][]uint32

	// 不同字节长度的ip信息,ip差值,内容位置
	diffLengthIpInfos [][]byte

	// 保留ip段
	reservedIpRanges []*ReservedIpRange

	contentArray []string

	// 内容位置字节长度
	contentIndexByteLength int

	licenseClient client.LicenseClient

	blockedIfRateLimited bool
}

func (client *IPv6GeoClient) Load(ctx *context.FastIPGeoContext) bool {
	geoConf := ctx.GeoConf
	client.blockedIfRateLimited = geoConf.BlockedIfRateLimited

	data := ctx.Data
	// 唯一性的内容条数
	contentArray := make([]string, xprt.ReadInt(data, consts.MetaInfoByteLength))
	client.contentArray = contentArray

	// 内容位置占字节数, 内容条数一般3字节足够了
	var contentIndexByteLength int
	if len(contentArray) < 1<<16 {
		contentIndexByteLength = 2
	} else {
		contentIndexByteLength = 3
	}
	client.contentIndexByteLength = contentIndexByteLength

	// 有多少种不同长度就生成多少个字节数组
	diffLengthIpInfos := make([][]byte, 16)
	client.diffLengthIpInfos = diffLengthIpInfos
	ipBlockStartIndex := make([][]uint32, 16)
	client.ipBlockStartIndex = ipBlockStartIndex
	ipBlockEndIndex := make([][]uint32, 16)
	client.ipBlockEndIndex = ipBlockEndIndex

	offset := consts.MetaInfoByteLength + 4 + 16*4

	for i := 0; i < 16; i++ {
		// 此字节长度下的ip数
		num := xprt.ReadInt(data, consts.MetaInfoByteLength+4+i*4)
		// 2字节的ip： 就存储内容序号
		// 6字节的ip ： 4字节后缀 + 4字节的ip差值 + 2/3字节的内容序号
		// 16字节的ip ： 14字节后缀 + 14字节的ip差值 + 2/3字节的内容序号
		if num > 0 {
			var length int
			if i < 2 {
				length = contentIndexByteLength * int(num)
			} else {
				length = (2*(i+1-2) + contentIndexByteLength) * int(num)
			}
			diffLengthIpInfos[i] = make([]byte, length)
			ipBlockStartIndex[i] = make([]uint32, Ipv6FirstSegmentSize)
			ipBlockEndIndex[i] = make([]uint32, Ipv6FirstSegmentSize)
		}
	}

	// 加载每个ip段的起始序号和结束序号
	for i := 0; i < 16; i++ {
		// 此长度的ip不存在
		if diffLengthIpInfos[i] == nil {
			offset += 8 * Ipv6FirstSegmentSize
			continue
		}
		for j := 0; j < Ipv6FirstSegmentSize; j++ {
			// 0:0:0:0:0:12:2e:3000  :  存储 6字节长度的ip 的 12 开头的ip, 他们的序号在 100~1000之间, 是6字节内部的序号
			ipBlockStartIndex[i][j] = xprt.ReadInt(data, offset)
			ipBlockEndIndex[i][j] = xprt.ReadInt(data, offset+4)
			offset += 8
		}
	}
	// 加载每个ip段的数据信息,ip后段, ip差值, 内容位置
	for i := 0; i < len(diffLengthIpInfos); i++ {
		bytes := diffLengthIpInfos[i]
		if bytes != nil {
			xprt.ArrayCopy(data, offset, bytes, 0, len(bytes))
			offset += len(bytes)
		}
	}

	// 阿里云用户id
	id := ctx.LicenseClient.GetId()
	version := ctx.MetaInfo.Version

	// 加载内容
	for i := 0; i < len(contentArray); i++ {
		length := int(xprt.ReadInt(data, offset))
		offset += 4
		// 将所有字符串都取出来, 每个字符串都缓存好
		rawContent := string(xprt.CopyOfRange(data, offset, offset+length))
		content := xprt.RawToJson(version, id, rawContent, ctx.MetaInfo.StoredProperties, geoConf.Properties)
		contentArray[i] = content
		offset += length
	}
	// 还有保留ip段
	if offset == len(data) {
		return true
	}
	reservedIpNum := int(xprt.ReadInt(data, offset))
	offset += 4
	reservedIpRanges := make([]*ReservedIpRange, reservedIpNum)
	client.reservedIpRanges = reservedIpRanges

	// 解析保留段
	for i := 0; i < reservedIpNum; i++ {
		length := xprt.ReadInt(data, offset)
		offset += 4
		rawContent := string(xprt.CopyOfRange(data, offset, offset+int(length)))
		splits := strings.SplitN(rawContent, ",", 3)
		index, _ := strconv.Atoi(splits[2])
		ipRange := ReservedIpRange{
			start:        toEffectiveByteArray(splits[0]),
			end:          toEffectiveByteArray(splits[1]),
			contentIndex: index,
		}
		reservedIpRanges[i] = &ipRange
	}
	// 对保留段做排序
	sort.Sort(ipRanges(reservedIpRanges))
	return true
}

func (client *IPv6GeoClient) Search(ip string) (string, error) {
	panic("implement me")
}

func toEffectiveByteArray(ipNum string) []byte {
	bigInt, _ := new(big.Int).SetString(ipNum, 10)
	bytes := bigInt.Bytes()
	length := calculateEffectiveLength(bytes)
	data := make([]byte, length)
	xprt.ArrayCopy(bytes, len(bytes)-length, data, 0, length)
	return data
}

func calculateEffectiveLength(bytes []byte) int {
	for i := 0; i < len(bytes); i++ {
		if bytes[i] != 0 {
			return len(bytes) - 1
		}
	}
	return 0
}

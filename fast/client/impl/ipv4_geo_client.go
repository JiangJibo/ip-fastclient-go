package ip_geo_client

import (
	"errors"
	"ip-fastclient-go/fast/consts"
	"ip-fastclient-go/fast/context"
	"ip-fastclient-go/fast/xprt"
	lsnClient "ip-fastclient-go/license/client"
	error2 "ip-fastclient-go/license/error"
	"strings"
)

var (
	IpFirstSegmentSize = 256 * 256
)

type Ipv4GeoClient struct {
	ipBlockStart         []uint32
	ipBlockEnd           []uint32
	endIpBytes           []byte
	contentIndexes       []byte
	contentArray         []string
	licenseClient        *lsnClient.LicenseClient
	blockedIfRateLimited bool
}

func (client *Ipv4GeoClient) Load(ctx *context.FastIPGeoContext) (bool, error) {
	client.licenseClient = ctx.LicenseClient
	client.blockedIfRateLimited = ctx.GeoConf.BlockedIfRateLimited
	geoConf := ctx.GeoConf
	data := ctx.Data

	ipBlockStart := make([]uint32, IpFirstSegmentSize)
	client.ipBlockStart = ipBlockStart

	ipBlockEnd := make([]uint32, IpFirstSegmentSize)
	client.ipBlockEnd = ipBlockEnd

	for i := 0; i < IpFirstSegmentSize; i++ {
		k := consts.MetaInfoByteLength + 8 + (i * 8)
		ipBlockStart[i] = xprt.ReadInt(data, k)
		ipBlockEnd[i] = xprt.ReadInt(data, k+4)
	}
	// 前4字节存储ip条数
	recordSize := xprt.ReadInt(data, consts.MetaInfoByteLength)

	endIpBytes := make([]byte, recordSize<<1)
	client.endIpBytes = endIpBytes
	contentIndexes := make([]byte, recordSize<<3)
	client.contentIndexes = contentIndexes
	// 有多少条唯一性的内容
	contentArray := make([]string, xprt.ReadInt(data, consts.MetaInfoByteLength+4))
	client.contentArray = contentArray
	// int形式的内容位置
	contentIndex := make([]int, recordSize)

	index := 0
	// 原始内容与处理过的内容间的映射
	contentMappings := make(map[string]string, 0)
	contentIndexMappings := make(map[string]int, 0)

	// 阿里云用户id
	id := ctx.LicenseClient.GetId()
	version := ctx.MetaInfo.Version
	metaInfo := ctx.MetaInfo

	for i := 0; i < int(recordSize); i++ {
		pos := consts.MetaInfoByteLength + 8 + IpFirstSegmentSize*8 + (i * 9)
		endIpBytes[2*i] = data[pos]
		endIpBytes[2*i+1] = data[pos+1]
		offset := xprt.ReadInt(data, pos+2)
		length := xprt.ReadVInt3(data, pos+6)
		if offset == 0 && length == 0 {
			continue
		}

		// 将所有字符串都取出来, 每个字符串都缓存好
		rawContent := string(xprt.CopyOfRange(data, int(offset), int(offset+length)))
		// 对原始内容做处理, 不重复处理
		var content string
		if v, ok := contentMappings[rawContent]; ok {
			content = v
		} else {
			content = xprt.RawToJson(version, id, rawContent, metaInfo.StoredProperties, geoConf.Properties)
			contentMappings[rawContent] = content
		}

		// 缓存字符串, 如果是一个新的字符串
		if _, ok := contentIndexMappings[content]; !ok {
			//因为我们为了在打包的时候插入一些水印ip，预先增加了20个ip段。一个水印ip一个单独的段，由于给的水印ip不一定在已有的两个段之间
			//这导致了ip的个数可能没有那么多，有些内容没有，存在超出index的情况，所以这里判断下
			if index > len(contentArray)-1 {
				continue
			}
			contentArray[index] = content
			contentIndexMappings[content] = index
			contentIndex[i] = index
			index++
		}
	}
	// 将内容位置的int数组转换成字节数组,节省一个字节
	for i := 0; i < len(contentIndex); i++ {
		xprt.WriteVInt3(contentIndexes, 3*i, contentIndex[i])
	}
	return true, nil
}

func (client *Ipv4GeoClient) Search(ip string) (string, error) {
	// TODO 限流， 这一步不能抽离出来，用来做代码加固
	if client.blockedIfRateLimited {
		client.licenseClient.Acquire()
	} else {
		b, err := client.licenseClient.TryAcquire()
		if err != error2.SUCCESS {
			return "", errors.New(err.Error())
		}
		if !b {
			return "", errors.New("be rate limited")
		}
	}

	// 计算ip前缀的int值
	firstDotIndex := strings.Index(ip, ".")
	firstSegmentInt := calculateIpSegmentInt(ip, 0, firstDotIndex-1)

	// 计算ip第二段的int值
	secondDotIndex := indexOf(ip, '.', firstDotIndex+1)
	secondSegmentInt := calculateIpSegmentInt(ip, firstDotIndex+1, secondDotIndex-1)
	prefixSegmentsInt := (firstSegmentInt << 8) + secondSegmentInt

	start := client.ipBlockStart[prefixSegmentsInt]
	end := client.ipBlockEnd[prefixSegmentsInt]
	suffix := calculateIpInteger(ip, 0, secondDotIndex+1)

	var cur int
	if start == end {
		cur = int(end)
	} else {
		cur = client.binarySearch(int(start), int(end), suffix)
	}
	return client.contentArray[xprt.ReadVInt3(client.contentIndexes, (cur<<1)+cur)], nil
}

//计算IP段的int值
func calculateIpSegmentInt(ip string, startIndex int, endIndex int) int {
	num := 0
	for i := startIndex; i <= endIndex; i++ {
		var radix int
		if i == endIndex {
			radix = 1
		} else if i == endIndex-1 {
			radix = 10
		} else {
			radix = 100
		}
		num = num + radix*(int(ip[i]-48))
	}
	return num
}

// 计算ip的int值
func calculateIpInteger(ip string, result int, dot int) int {
	var num int
	for {
		dotIndex := indexOf(ip, '.', dot) - 1
		if dotIndex >= 0 {
			num = calculateIpSegmentInt(ip, dot, dotIndex)
		} else {
			num = calculateIpSegmentInt(ip, strings.LastIndex(ip, ".")+1, len(ip)-1)
		}
		result <<= 8
		result |= num & 0xff
		dot = dotIndex + 2
		if dotIndex < 0 {
			return result
		}
	}
}

// 从某个fromIndex位置开始寻找text第一个出现dot的下标
func indexOf(text string, dot uint8, fromIndex int) int {
	for i := fromIndex; i < len(text); i++ {
		if text[i] == dot {
			return i
		}
	}
	return 0
}

func (client *Ipv4GeoClient) binarySearch(low int, high int, suffix int) int {
	mid := 0
	for {
		if low > high {
			break
		}
		mid = (low + high) >> 1
		switch client.compareSuffixBytes(mid, suffix) {
		case 1:
			high = mid
			break
		case 0:
			return mid
		case -1:
			low = mid + 1
		}
	}
	return mid
}

// 对比IP的后两个字节
func (client *Ipv4GeoClient) compareSuffixBytes(index int, suffix int) int {
	value := int(client.endIpBytes[index<<1]&0xff)<<8 | int(client.endIpBytes[(index<<1)+1])&0xff
	if value > suffix {
		return 1
	} else if value == suffix {
		return 0
	} else {
		return -1
	}
}

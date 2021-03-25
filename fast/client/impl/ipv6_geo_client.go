package impl

import (
	"errors"
	"fmt"
	"github.com/jiangjibo/ip-fastclient-go/fast/consts"
	"github.com/jiangjibo/ip-fastclient-go/fast/context"
	"github.com/jiangjibo/ip-fastclient-go/fast/xprt"
	"github.com/jiangjibo/ip-fastclient-go/license/client"
	lsError "github.com/jiangjibo/ip-fastclient-go/license/error"
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

	licenseClient *client.LicenseClient

	blockedIfRateLimited bool
}

func (client *IPv6GeoClient) Load(ctx *context.FastIPGeoContext) bool {
	geoConf := ctx.GeoConf
	client.blockedIfRateLimited = geoConf.BlockedIfRateLimited
	client.licenseClient = ctx.LicenseClient

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
	// 被限流时阻塞
	if client.blockedIfRateLimited {
		_, err := client.licenseClient.Acquire()
		if err != lsError.SUCCESS {
			return "", errors.New(err.Error())
		}
	} else {
		b, err := client.licenseClient.TryAcquire()
		// 有license异常
		if err != lsError.SUCCESS {
			return "", errors.New(err.Error())
		}
		// 被限流
		if !b {
			return "", errors.New("rate limited")
		}
	}
	// 将ip转换成字节数组
	ipv6Address, err := ToByteArray(ip)
	if err != nil {
		return "", err
	}
	// 当前字节长度在所有长度中的序号, 从0开始
	segmentIndex := len(ipv6Address) - 1

	// 当前长度下的ip没有
	if segmentIndex < 0 || client.ipBlockStartIndex[segmentIndex] == nil {
		return client.searchInReservedIpRanges(ipv6Address), nil
	}

	var prefixSegmentsInt int
	if len(ipv6Address) == 1 { // 一个字节的ip
		prefixSegmentsInt = int(ipv6Address[0])
	} else {
		prefixSegmentsInt = int(ipv6Address[0])<<8 + int(ipv6Address[1])
	}

	start := client.ipBlockStartIndex[segmentIndex][prefixSegmentsInt]
	end := client.ipBlockEndIndex[segmentIndex][prefixSegmentsInt]

	if start < 0 {
		return client.searchInReservedIpRanges(ipv6Address), nil
	}
	index := client.binarySearch(int(start), int(end), ipv6Address)
	if index > 0 {
		return client.contentArray[index], nil
	}
	return client.searchInReservedIpRanges(ipv6Address), nil
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

func ToByteArray(address string) ([]byte, error) {
	if address == "" {
		return nil, errors.New(fmt.Sprintf("Invalid length - the string %s is too short to be an IPv6 address", address))
	}
	// 验证长度
	first := 0
	last := len(address)
	if last < 2 {
		return nil, errors.New(fmt.Sprintf("Invalid length - the string %s is too short to be an IPv6 address", address))
	}
	length := last - first
	if length > 39 { // 32个数字 + 7个":"
		return nil, errors.New(fmt.Sprintf("Invalid length - the string %s is too long to be an IPv6 address. Length: %d", address, last))
	}

	partIndex := 0
	partHexDigitCount := 0
	afterDoubleSemicolonIndex := last + 2

	var data []byte
	var v byte
	k := 0
	s := 0
	x := 1
	left := 0
	for i := first; i < last; i++ {
		c := address[i]
		if isHexDigit(c) {
			if c == '0' && data == nil {
				continue
			}
			// 定位当前段有几个数字, 且当前数字的序号：1 2 3 4 中的一个
			if partHexDigitCount == 0 {
				y := i + 1
				for i := 0; i < 4; i++ {
					y++
					if y-1 >= last || address[y-1] == ':' {
						break
					} else {
						x++
					}
				}
				// ipv6分8段, 每一段都会写入两个字节
				left += 2
				// 当前数字的序号
				partHexDigitCount = 4 - x
			}
			// 初始化字节数组
			if data == nil {
				var l int
				if x > 2 {
					l = 2
				} else {
					l = 1
				}
				dataLength := ((7 - partIndex) << 1) + l
				data = make([]byte, dataLength)
			}
			// 每段数字不能超过4个
			partHexDigitCount++
			if partHexDigitCount > 4 {
				return nil, errors.New(fmt.Sprintf("Ipv6 address %s parts must contain no more than 16 bits (4 hex digits)", address))
			}
			// 当前数字的字节值
			if c >= 97 {
				v = c - 87
			} else {
				v = c - 48
			}
			// 0:0:12:c::12:226 将'c'的前一字节填充0
			if x == 1 && partHexDigitCount == 4 && k > 0 {
				data[k] = 0
				s++
				k++
			}
			// 0:0:12:2c::12:226 将'2c'的前一字节填充0
			if x == 2 && partHexDigitCount == 3 && k > 0 {
				data[k] = 0
				s++
				k++
			}
			// 根据数字位置填充字节的前一半
			if partHexDigitCount == 1 || partHexDigitCount == 3 {
				data[k] = v << 4
			}
			// 填充字节的后一半
			if partHexDigitCount == 2 || partHexDigitCount == 4 {
				// 如果此段只有一个数字,或者有3个数字，当前是第一个数字
				if x == 1 || (partHexDigitCount == 2 && x == 3) {
					data[k] = 0
				}
				data[k] = data[k] | v
				s++
				k++
			}
		} else {
			if c == ':' {
				if data != nil {
					// 下一次写的字节位置
					if k != 1 {
						k += 2 - s
					}
				}
				s = 0
				x = 1
				partIndex++
				partHexDigitCount = 0
				// 如果存在连续的两个":", 即 "::"
				if i < last-1 && address[i+1] == ':' {
					// 在两个: 之后的下一个数字的位置
					afterDoubleSemicolonIndex = i + 2
					break
				}
				continue
			}
			return nil, errors.New(fmt.Sprintf("Ipv6 address %s illegal character: %c at index %d", address, c, i))
		}
	}

	if data == nil {
		data = make([]byte, 16)
	}

	// 从末尾倒叙向前遍历 直至 ::
	lastFilledPartIndex := partIndex - 1
	l := len(data) - 1
	right := l
	markRight := l
	partIndex = 7

	for i := last - 1; i >= afterDoubleSemicolonIndex; i-- {
		c := address[i]

		if isHexDigit(c) {
			if partIndex <= lastFilledPartIndex {
				return nil, errors.New(fmt.Sprintf("Ipv6 address %s too many parts. Expected 8 parts", address))
			}
			partHexDigitCount++

			if partHexDigitCount > 4 {
				return nil, errors.New(fmt.Sprintf("Ipv6 address %s parts must contain no more than 16 bits (4 hex digits)", address))
			}
			// 当前数字的字节值
			if c >= 97 {
				v = c - 87
			} else {
				v = c - 48
			}
			// 根据数字位置填充字节的后一半
			if partHexDigitCount == 1 || partHexDigitCount == 3 {
				right--
				data[l] = v
			}
			// 填充字节的前一半
			if partHexDigitCount == 2 || partHexDigitCount == 4 {
				data[l] = data[l] | v<<4
				s++
				l--
			}
			if c != '0' {
				markRight = right
			}
		} else {
			if c == ':' {
				if partHexDigitCount < 3 {
					right--
				}
				l -= 2 - s
				s = 0
				partIndex--
				partHexDigitCount = 0
				continue
			}
			return nil, errors.New(fmt.Sprintf("Ipv6 address %s illegal character: %c at index %d", address, c, i))
		}
	}

	// 0::X... 类型的IP
	if left == 0 {
		poolOffset := 16 - markRight - 1
		if poolOffset <= 0 {
			return nil, nil
		}
		tmp := make([]byte, poolOffset)
		for i := 0; i < len(tmp); i++ {
			tmp[i] = data[markRight+i+1]
		}
		data = tmp
	} else {
		// 填充 "::" 代表的空字节, 置为0, 因为字节数组是缓存重用的, 需要复位
		if left != right {
			for i := left; i <= right; i++ {
				data[i] = 0
			}
		}
	}
	return data, nil
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func (client IPv6GeoClient) searchInReservedIpRanges(iPv6Address []byte) string {
	index := client.searchReservedIpRanges(iPv6Address)
	if index == -1 {
		return ""
	} else {
		return client.contentArray[index]
	}
}

// 检索保留ip段
func (client IPv6GeoClient) searchReservedIpRanges(iPv6Address []byte) int {
	if client.reservedIpRanges == nil || len(client.reservedIpRanges) == 0 {
		return -1
	}
	length := len(iPv6Address)
	for _, ipRange := range client.reservedIpRanges {
		l1 := len(ipRange.start)
		l2 := len(ipRange.end)
		if l1 > length || l2 < length {
			continue
		}
		if l1 < length && l2 > length {
			return ipRange.contentIndex
		}
		// 长度和起始ip段相等
		if l1 == length {
			// 如果是0
			if l1 == 0 {
				return ipRange.contentIndex
			}
			for i := 0; i < length; i++ {
				if ipRange.start[i] < (iPv6Address[i]) {
					return ipRange.contentIndex
				}
				// 搜索ip是保留段的起始ip
				if i == length-1 && (ipRange.start[i]) == (iPv6Address[i]) {
					return ipRange.contentIndex
				}
			}
		}
		// 长度和结束ip段相等
		if l2 == length {
			for i := 0; i < length; i++ {
				if (ipRange.end[i]) > (iPv6Address[i]) {
					return ipRange.contentIndex
				}
				// 搜索ip是保留段的结束ip
				if i == length-1 && (ipRange.end[i]) == (iPv6Address[i]) {
					return ipRange.contentIndex
				}
			}
		}
	}
	return -1
}

/**
 * 定位ip序号
 *
 * @param low  ip段的起始序号
 * @param high ip段的结束序号
 * @param ip   IP字节数组, 16字节
 * @return int
 */
func (client IPv6GeoClient) binarySearch(low, high int, ip []byte) int {
	// 存储当前长度的ip信息的数组
	data := client.diffLengthIpInfos[len(ip)-1]

	// 如果是一字节或者二字节的ip
	if len(ip) <= 2 {
		return int(readVint(data, high*client.contentIndexByteLength, client.contentIndexByteLength))
	}
	// 每份ip信息的字节长度
	perLength := ((len(ip) - 2) << 1) + client.contentIndexByteLength

	result := -2
	order := high

	for low <= high {
		mid := (low + high) >> 1
		result = compareByteArray(data, mid*perLength, ip, 2)
		// 如果结束的ip就是当前ip
		if result == 0 {
			order = mid
			break
		} else if result > 0 {
			order = mid
			if mid == 0 {
				break
			}
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	offset := order * perLength
	if result != 0 && !compareByStartIp(ip, 2, data, offset+len(ip)-2) {
		return -1
	}
	return int(readVint(data, offset+((len(ip)-2)<<1), client.contentIndexByteLength))
}

// 读取3个字节凑成Int
func readVint(data []byte, p, length int) uint32 {
	if length == 2 {
		x := uint32(data[p]) << 8 & 0xFF00
		p++
		y := uint32(data[p])
		return x | y
	} else {
		x := uint32(data[p]) << 16 & 0xFF0000
		p++
		y := uint32(data[p]) << 8 & 0xFF00
		p++
		z := uint32(data[p])
		return x | y | z
	}
}

/**
 * 比较两个字节数字的大小,逐位比较
 *
 * @param array1 数组1
 * @param index1 开始的比较位
 * @param array2 数组2
 * @param index2 开始的比较位
 * @return int
 */
func compareByteArray(array1 []byte, index1 int, array2 []byte, index2 int) int {
	k := 0
	result := 0
	for i := index1; i < index1+len(array2); i++ {
		if array1[i] != array2[index2+k] {
			return int(array1[i]) - int(array2[index2+k])
		}
		// end ip 和当前检索ip相同
		k++
		if index2+k == len(array2) {
			return 0
		}
	}
	return result
}

/**
 * 和起始ip比较
 *
 * @param ip         检索ip
 * @param index
 * @param data       startIp后缀
 * @param startIndex
 * @return
 */
func compareByStartIp(ip []byte, index int, data []byte, startIndex int) bool {
	k := 0
	for i := index; i < len(ip); i++ {
		x1 := ip[i]
		x2 := data[startIndex+k]
		if x1 < x2 {
			return false
		} else if x1 > x2 {
			return true
		}
		k++
	}
	return true
}

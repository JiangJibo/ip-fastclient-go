package ip_geo_client

import (
	"ip-fastclient-go/fast/consts"
	"ip-fastclient-go/fast/context"
	"ip-fastclient-go/fast/xprt"
	lsnClient "ip-fastclient-go/license/client"
)

var (
	IpFirstSegmentSize = 256 * 256
)

type Ipv4GeoClient struct {
	ipBlockStart   []uint32
	ipBlockEnd     []uint32
	endIpBytes     []byte
	contentIndexes []byte
	contentArray   []string
	licenseClient  lsnClient.LicenseClient
}

func (client *Ipv4GeoClient) Search(ip string) (string, error) {
	return "", nil
}

func (client *Ipv4GeoClient) Load(ctx context.FastIPGeoContext) (bool, error) {
	geoConf := ctx.GeoConf
	data := ctx.Data

	ipBlockStart := make([]uint32, IpFirstSegmentSize)
	client.ipBlockStart = ipBlockStart

	ipBlockEnd := make([]uint32, IpFirstSegmentSize)
	client.ipBlockEnd = ipBlockEnd

	for i := 0; i < IpFirstSegmentSize; i++ {
		k := consts.META_INFO_BYTE_LENGTH + 8 + (i * 8)
		ipBlockStart[i] = clientUtils.ReadInt(data, k)
		ipBlockEnd[i] = clientUtils.ReadInt(data, k+4)
	}
	// 前4字节存储ip条数
	recordSize := clientUtils.ReadInt(data, consts.META_INFO_BYTE_LENGTH)

	endIpBytes := make([]byte, recordSize<<1)
	client.endIpBytes = endIpBytes
	contentIndexes := make([]byte, recordSize<<3)
	client.contentIndexes = contentIndexes
	// 有多少条唯一性的内容
	contentArray := make([]string, clientUtils.ReadInt(data, consts.META_INFO_BYTE_LENGTH+4))
	client.contentArray = contentArray
	// int形式的内容位置
	contentIndex := make([]uint32, recordSize)

	index := 0
	// 原始内容与处理过的内容间的映射
	contentMappings := make(map[string]string, 0)
	contentIndexMappings := make(map[string]uint32, 0)

	for i := 0; i < int(recordSize); i++ {
		pos := consts.META_INFO_BYTE_LENGTH + 8 + IpFirstSegmentSize + (i * 9)
		endIpBytes[2*i] = data[pos]
		endIpBytes[2*i+1] = data[pos+1]
		offset := clientUtils.ReadInt(data, pos+2)
		length := clientUtils.ReadInt(data, pos+6)
		if offset == 0 && length == 0 {
			continue
		}
		// 将所有字符串都取出来, 每个字符串都缓存好
		rawContent := string(clientUtils.CopyOfRange(data, int(offset), int(offset+length)))
		// 对原始内容做处理, 不重复处理
		var content string
		if v, ok := contentMappings[rawContent]; ok {
			content = v
		} else {

		}

	}

	// 阿里云用户id

	return true, nil
}

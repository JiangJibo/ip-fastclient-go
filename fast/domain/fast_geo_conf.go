package domain

import (
	"github.com/jiangjibo/ip-fastclient-go/fast/xprt"
	"io"
	"io/ioutil"
	"os"
)

type FastGeoConf struct {
	Properties map[string]bool

	LicenseFilePath string
	LicenseBytes    []byte
	LicenseInput    io.Reader

	DataFilePath string
	DataBytes    []byte
	DataInput    io.Reader

	BlockedIfRateLimited bool
}

func (fastGeoConf *FastGeoConf) FilterEmptyValue() {
	os.Setenv(xprt.IpSdkFilterEmptyPropertyConfKey, "true")
}

// 读取Reader
func ReadInput(reader io.Reader) []byte {
	if reader == nil {
		return nil
	}
	data, _ := ioutil.ReadAll(reader)
	return data
}

// 获取许可数据
func (fastGeoConf *FastGeoConf) GetLicenseData() []byte {
	if fastGeoConf.LicenseInput != nil {
		fastGeoConf.initLicense()
		return fastGeoConf.LicenseBytes
	}
	if fastGeoConf.LicenseBytes != nil {
		return fastGeoConf.LicenseBytes
	}
	bts, _ := ioutil.ReadFile(fastGeoConf.LicenseFilePath)
	return bts
}

// 初始化许可数据
func (fastGeoConf *FastGeoConf) initLicense() error {
	if fastGeoConf.LicenseInput == nil {
		return nil
	}
	licenseBytes := ReadInput(fastGeoConf.LicenseInput)
	if licenseBytes != nil || len(licenseBytes) > 0 {
		fastGeoConf.LicenseBytes = licenseBytes
		fastGeoConf.LicenseInput = nil
	}
	return nil
}

// 获取数据
func (fastGeoConf *FastGeoConf) GetDexData() []byte {
	if fastGeoConf.DataInput != nil {
		fastGeoConf.initDex()
		return fastGeoConf.DataBytes
	}
	if fastGeoConf.DataBytes != nil {
		return fastGeoConf.DataBytes
	}
	bts, _ := ioutil.ReadFile(fastGeoConf.DataFilePath)
	return bts
}

// 初始化数据
func (fastGeoConf *FastGeoConf) initDex() error {
	if fastGeoConf.DataInput == nil {
		return nil
	}
	dataBytes := ReadInput(fastGeoConf.DataInput)
	if dataBytes != nil || len(dataBytes) > 0 {
		fastGeoConf.DataBytes = dataBytes
		fastGeoConf.DataInput = nil
	}
	return nil
}

func (fastGeoConf *FastGeoConf) CalculateLicenseKey() string {
	if fastGeoConf.LicenseInput != nil {
		return xprt.Md5Hex(ReadInput(fastGeoConf.LicenseInput))
	}
	if fastGeoConf.LicenseBytes != nil {
		return xprt.Md5Hex(fastGeoConf.LicenseBytes)
	}
	return fastGeoConf.LicenseFilePath
}

// 释放资源
func (fastGeoConf *FastGeoConf) ReleaseResources() {
	if fastGeoConf.DataInput != nil {
		fastGeoConf.DataInput = nil
	}
	if fastGeoConf.DataBytes != nil {
		fastGeoConf.DataBytes = nil
	}
	if fastGeoConf.LicenseInput != nil {
		fastGeoConf.LicenseInput = nil
	}
	if fastGeoConf.LicenseBytes != nil {
		fastGeoConf.LicenseBytes = nil
	}
}

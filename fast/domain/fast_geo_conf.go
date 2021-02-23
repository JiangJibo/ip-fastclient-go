package domain

import (
	"io"
	"io/ioutil"
	"os"
)

var (
	IpSdkFilterEmptyPropertyConfKey = "ip.sdk.filter.empty.property"
)

type FastGeoConf struct {
	Properties map[string]bool

	LicenseFilePath string
	LicenseBytes    []byte
	LicenseInput    io.Reader

	DataFilePath string
	DataBytes    []byte
	DataInput    io.Reader
}

func (fastGeoConf *FastGeoConf) FilterEmptyValue() {
	os.Setenv(IpSdkFilterEmptyPropertyConfKey, "true")
}

// 读取Reader
func ReadInput(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, nil
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 获取许可数据
func (fastGeoConf *FastGeoConf) GetLicenseData() ([]byte, error) {
	if fastGeoConf.LicenseInput != nil {
		fastGeoConf.initLicense()
		return fastGeoConf.LicenseBytes, nil
	}
	if fastGeoConf.LicenseBytes != nil {
		return fastGeoConf.LicenseBytes, nil
	}
	return ioutil.ReadFile(fastGeoConf.LicenseFilePath)
}

// 初始化许可数据
func (fastGeoConf *FastGeoConf) initLicense() error {
	if fastGeoConf.LicenseInput == nil {
		return nil
	}
	licenseBytes, err := ReadInput(fastGeoConf.LicenseInput)
	if err != nil {
		return err
	}
	if licenseBytes != nil || len(licenseBytes) > 0 {
		fastGeoConf.LicenseBytes = licenseBytes
		fastGeoConf.LicenseInput = nil
	}
	return nil
}

// 获取数据
func (fastGeoConf *FastGeoConf) GetDexData() ([]byte, error) {
	if fastGeoConf.DataInput != nil {
		fastGeoConf.initDex()
		return fastGeoConf.DataBytes, nil
	}
	if fastGeoConf.DataBytes != nil {
		return fastGeoConf.DataBytes, nil
	}
	return ioutil.ReadFile(fastGeoConf.DataFilePath)
}

// 初始化数据
func (fastGeoConf *FastGeoConf) initDex() error {
	if fastGeoConf.DataInput == nil {
		return nil
	}
	dataBytes, err := ReadInput(fastGeoConf.DataInput)
	if err != nil {
		return err
	}
	if dataBytes != nil || len(dataBytes) > 0 {
		fastGeoConf.DataBytes = dataBytes
		fastGeoConf.DataInput = nil
	}
	return nil
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

package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	LicenseConsts "github.com/JiangJibo/ip-fastclient-go/license/consts"
	LicenseDomain "github.com/JiangJibo/ip-fastclient-go/license/domain"
	LicenseErrors "github.com/JiangJibo/ip-fastclient-go/license/error"
	LicenseUtils "github.com/JiangJibo/ip-fastclient-go/license/utils"
	"golang.org/x/time/rate"
	"io/ioutil"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	mutex            = new(sync.Mutex)
	licenseClientMap = make(map[string]*LicenseClient, 2)
)

type LicenseClient struct {

	// 证书数据
	LicenseData []byte

	// 证书数据文件地址
	LicenseFilePath string

	// 限流
	rateLimiter *rate.Limiter

	//是否不限速
	isUnlimited bool

	//原子计数，用来判断证书是否过期或者使用qps是否超过售卖限制
	checker atomic.Value
}

//单例返回licenseClient, 考虑到用户可能同时使用ipv4/ipv6, 或者多个v4/v6, 以为license路径作为单例主键
func GetInstance(licenseKey string, licenseData []byte) *LicenseClient {
	return getInstance(licenseKey, licenseData)
}

//单例返回licenseClient, 考虑到用户可能同时使用ipv4/ipv6, 或者多个v4/v6, 以为license路径作为单例主键
func GetInstanceByLicensePath(licenseFilePath string) *LicenseClient {
	return getInstance(licenseFilePath, nil)
}

func getInstance(licenseKey string, licenseData []byte) *LicenseClient {
	if _, ok := licenseClientMap[licenseKey]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if _, ok := licenseClientMap[licenseKey]; !ok {
			licenseClient := LicenseClient{
				LicenseData:     licenseData,
				LicenseFilePath: licenseKey,
			}
			licenseClient.Init()
			licenseClientMap[licenseKey] = &licenseClient
		}

	}
	return licenseClientMap[licenseKey]
}

func (lc *LicenseClient) Init() LicenseErrors.LicenseError {
	if lc.LicenseData == nil {
		if lc.LicenseFilePath == "" {
			return LicenseErrors.LicenseFileNotExists
		}
		data, _ := ioutil.ReadFile(lc.LicenseFilePath)
		lc.LicenseData = data
	}
	return lc.doInit()
}

func (lc *LicenseClient) TryAcquire() (bool, LicenseErrors.LicenseError) {
	v := lc.checker.Load()
	if v != nil {
		return false, v.(LicenseErrors.LicenseError)
	}
	if lc.isUnlimited {
		return true, LicenseErrors.SUCCESS
	}
	return lc.rateLimiter.Allow(), LicenseErrors.SUCCESS
}

func (lc *LicenseClient) Acquire() (bool, LicenseErrors.LicenseError) {
	v := lc.checker.Load()
	if v != nil {
		return false, v.(LicenseErrors.LicenseError)
	}
	if lc.isUnlimited {
		return true, LicenseErrors.SUCCESS
	}
	allow := lc.rateLimiter.Allow()
	if !allow {
		lc.rateLimiter.Wait(context.TODO())
	}
	return true, LicenseErrors.SUCCESS
}

func (lc *LicenseClient) XTry() string {
	rnd := LicenseUtils.CreateRandomNumber(32)
	secret := lc.decryptLicense()
	word := secret.IsValidate()
	isValid := LicenseUtils.Decox(word)
	if isValid {
		return LicenseUtils.Echo("")
	} else {
		return rnd
	}
}

func (lc *LicenseClient) GetDataType() string {
	return lc.decryptLicense().GetDataType()
}

func (lc *LicenseClient) GetId() string {
	licenseSecret := lc.decryptLicense()
	id, lsErr := licenseSecret.GetId()
	if lsErr != LicenseErrors.SUCCESS {
		panic(lsErr.Error())
	}
	return id
}

func (lc *LicenseClient) doInit() LicenseErrors.LicenseError {
	//第一次初始化
	licenseSecret, err := lc.firstInit()
	if err != LicenseErrors.SUCCESS {
		return err
	}
	//RateLimiter初始化
	lsnRateLimit := licenseSecret.GetRateLimit()
	lc.isUnlimited = LicenseConsts.MgUnlimitedQps == lsnRateLimit
	if !lc.isUnlimited {
		rl, err := strconv.Atoi(lsnRateLimit)
		if err != nil {
			return LicenseErrors.LicenseInvalid
		}
		// 每秒产生rl个令牌， 最多存储rl个令牌
		lc.rateLimiter = rate.NewLimiter(rate.Limit(rl), rl)
	}
	lc.lsnCheckerInit()
	return LicenseErrors.SUCCESS
}

func (lc *LicenseClient) firstInit() (*LicenseDomain.LicenseSecret, LicenseErrors.LicenseError) {
	licenseSecret := lc.decryptLicense()
	word := licenseSecret.IsValidate()
	bool := LicenseUtils.Decox(word)
	if !bool {
		return licenseSecret, LicenseErrors.LicenseInvalid
	}
	return licenseSecret, LicenseErrors.SUCCESS
}

func (lc *LicenseClient) decryptLicense() *LicenseDomain.LicenseSecret {

	fileContentPlainBytes := make([]byte, len(lc.LicenseData))
	_, err := base64.StdEncoding.Decode(fileContentPlainBytes, lc.LicenseData)
	if err != nil {
		panic(err)
	}
	fileContent := string(fileContentPlainBytes)

	//aes解密
	aesPassword := getAesPassword(fileContent)
	encryptContent := getEncryptContent(fileContent)
	decryptContent := LicenseUtils.AesDecryptECB(encryptContent, []byte(aesPassword))

	//rsa解密
	var license LicenseDomain.License
	err = json.Unmarshal(decryptContent, &license)
	if err != nil {
		panic(err)
	}
	cipherBytes := license.CipherBytes
	plainPublicKey := getClearlyPbk(license.PublicKey)
	cipherEntityDecryptBytes := LicenseUtils.DecryptByPublicKey(plainPublicKey, cipherBytes)

	var cipherEntity LicenseDomain.CipherEntity
	err = json.Unmarshal(cipherEntityDecryptBytes, &cipherEntity)
	if err != nil {
		panic(err)
	}
	ls := LicenseDomain.License{
		Sign:        license.Sign,
		PublicKey:   plainPublicKey,
		CipherBytes: cipherBytes,
	}
	return &LicenseDomain.LicenseSecret{
		License:      &ls,
		CipherEntity: &cipherEntity,
	}
}

// 每隔30分钟检查license文件
func (lc *LicenseClient) lsnCheckerInit() {
	for range time.Tick(time.Minute * 30) {
		err := lc.lsnCheck()
		if err != LicenseErrors.SUCCESS {
			lc.checker.Store(err)
		}
	}
}

func (lc *LicenseClient) lsnCheck() LicenseErrors.LicenseError {
	secret := lc.decryptLicense()
	word := secret.IsValidate()
	isValid := LicenseUtils.Decox(word)
	if !isValid {
		return LicenseErrors.UNKNOWN
	}
	return LicenseErrors.SUCCESS
}

func getAesPassword(fileContent string) string {
	return fileContent[0:LicenseConsts.MgAesPasswordSize]
}

func getEncryptContent(fileContent string) string {
	return fileContent[LicenseConsts.MgAesPasswordSize:]
}

func getClearlyPbk(pbk string) string {
	totalPrefixLength := LicenseConsts.MgPrefixConfusionSize + 1
	totalPostfixLength := LicenseConsts.MgPostfixConfusionSize + 1
	return pbk[totalPrefixLength : len(pbk)-totalPostfixLength]
}

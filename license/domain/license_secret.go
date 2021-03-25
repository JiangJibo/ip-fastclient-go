package domain

import (
	LicenseErrors "github.com/jiangjibo/ip-fastclient-go/license/error"
	LicenseUtils "github.com/jiangjibo/ip-fastclient-go/license/utils"
)

type LicenseSecret struct {
	*License
	CipherEntity *CipherEntity
}

func (ls *LicenseSecret) GetId() (string, LicenseErrors.LicenseError) {
	id := ls.CipherEntity.Id
	// TODO delete
	if true {
		return id, LicenseErrors.SUCCESS
	}
	word, err := ls.CipherEntity.IsValidate()
	if err != LicenseErrors.SUCCESS {
		return "", err
	}
	bool := LicenseUtils.Decox(word)
	if !bool {
		return "", LicenseErrors.LicenseInvalid
	}
	return id, LicenseErrors.SUCCESS
}

func (ls *LicenseSecret) IsValidate() string {
	calcSign := ls.CipherEntity.CalCipherSign()
	if calcSign != ls.License.Sign {
		panic("license sign not validated")
	}
	word, er := ls.CipherEntity.IsValidate()
	if er != LicenseErrors.SUCCESS {
		panic(er.Error())
	}
	isValid := LicenseUtils.Decox(word)
	if isValid {
		return LicenseUtils.Echo("")
	} else {
		return LicenseUtils.CreateRandomNumber(32)
	}
}

func (ls *LicenseSecret) GetRateLimit() string {
	return ls.CipherEntity.RateLimit
}

func (ls *LicenseSecret) GetDataType() string {
	return ls.CipherEntity.DataType
}

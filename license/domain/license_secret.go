package domain

import (
	"errors"
	LicenseErrors "ip-fastclient-go/license/error"
	LicenseUtils "ip-fastclient-go/license/utils"
)

type LicenseSecret struct {
	*License
	CipherEntity *CipherEntity
}

func (ls *LicenseSecret) GetId() (string, LicenseErrors.LicenseError) {
	id := ls.CipherEntity.Id
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

func (ls *LicenseSecret) IsValidate() (string, error) {
	calcSign, err := ls.CipherEntity.CalCipherSign()
	if err != nil {
		return "", err
	}
	if calcSign != ls.License.Sign {
		return "", nil
	}
	word, er := ls.CipherEntity.IsValidate()
	if er != LicenseErrors.SUCCESS {
		return "", errors.New(er.Error())
	}
	isValid := LicenseUtils.Decox(word)
	if isValid {
		return LicenseUtils.Echo(""), nil
	} else {
		return LicenseUtils.CreateRandomNumber(32), nil
	}
}

func (ls *LicenseSecret) GetRateLimit() string {
	return ls.CipherEntity.RateLimit
}

func (ls *LicenseSecret) GetDataType() string {
	return ls.CipherEntity.DataType
}

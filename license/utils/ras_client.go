package LicenseUtils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	gorsa "ip-fastclient-go/license/rsa"
)

// 使用公钥解密
func DecryptByPublicKey(publicKey string, cipherBytes []byte) ([]byte, error) {
	pk, err := getPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	grsa := gorsa.RSASecurity{}
	grsa.SetPubKey(pk)
	return grsa.PubKeyDECRYPT(cipherBytes)
}

func getPublicKey(publicKey string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), err
}

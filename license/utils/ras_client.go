package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	gorsa "ip-fastclient-go/license/gorsa"
)

// 使用公钥解密
func DecryptByPublicKey(publicKey string, cipherBytes []byte) []byte {
	pk := getPublicKey(publicKey)
	grsa := gorsa.RSASecurity{}
	grsa.SetPubKey(pk)
	ret, err := grsa.PubKeyDECRYPT(cipherBytes)
	if err != nil {
		panic(err)
	}
	return ret
}

func getPublicKey(publicKey string) *rsa.PublicKey {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	pub, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		panic(err)
	}
	return pub.(*rsa.PublicKey)
}

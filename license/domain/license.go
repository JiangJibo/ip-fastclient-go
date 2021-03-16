package domain

type License struct {
	Sign      string
	PublicKey string
	//消息里面的加密部分
	CipherBytes []byte
}

package encrpt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github/wziww/medusa/log"
	"io"
)

// AesCfb ...
type AesCfb struct {
	Password    *[]byte
	PaddingMode string
}

var _ Encryptor = (*AesCfb)(nil)

// NewAesCfb constructor...
// func NewAesCfb(password *[]byte, iv *[]byte) *aesCfb {
// 	if len(*password) != 16 && len(*password) != 24 && len(*password) != 32 {
// 		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: password长度必须为16、24或32位"))
// 		return nil
// 	}
// 	if len(*iv) != 16 {
// 		log.FMTLog(log.LOGERROR, errors.New("aes_ctr: iv长度必须为16位"))
// 		return nil
// 	}
// 	ctr := &aesCfb{password, iv}
// 	return ctr
// }

// Decode ...
func (st *AesCfb) Decode(cipherBuf []byte) []byte {

	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	if len(cipherBuf) < aes.BlockSize {
		log.FMTLog(log.LOGERROR, errors.New("aes_cfb: ciphertext too short"))
		return nil
	}
	iv := cipherBuf[:aes.BlockSize]
	var buf = cipherBuf[aes.BlockSize:]
	// unpad
	buf, _ = HandleUnPadding(st.PaddingMode)(buf, aes.BlockSize)
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(buf, buf)

	return buf
}

// Encode ...
func (st *AesCfb) Encode(plainBuf []byte) []byte {
	block, err := aes.NewCipher(*st.Password)
	if err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}
	// pad
	plainBuf = HandlePadding(st.PaddingMode)(plainBuf, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plainBuf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.FMTLog(log.LOGERROR, err)
		return nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainBuf)
	return ciphertext
}

// Construct ...
func (st *AesCfb) Construct(name string) interface{} {
	var targetKeySize int
	switch name {
	case "aes-128-cfb":
		targetKeySize = 16
	case "aes-192-cfb":
		targetKeySize = 24
	case "aes-256-cfb":
		targetKeySize = 32
	default:
		return nil
	}
	if len(*st.Password) != targetKeySize {
		return nil
	}
	return st
}

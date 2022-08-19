package secret

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

func Sm4Encrypt(plainText string, key string, iv string) (string, error) {
	var out string
	// fmt.Printf("iv: %s", sm4.IV)
	sm4.SetIV([]byte(iv))
	sm4_out, err := sm4.Sm4Cbc([]byte(key), []byte(plainText), true)
	if err != nil {
		return out, err
	}
	out = hex.EncodeToString(sm4_out)
	return out, err
}

func Sm4Decrypt(cipherText string, key string, iv string) (string, error) {
	var out string
	cipherByte, err := hex.DecodeString(cipherText)
	if err != nil {
		return out, err
	}
	sm4.SetIV([]byte(iv))
	sm4_out, err := sm4.Sm4Cbc([]byte(key), cipherByte, false)
	if err != nil {
		return out, err
	}
	out = string(sm4_out)
	return out, err
}

// func Sm4DecryptWithKey(cipherText string, key string) (string, error) {
// 	var out string
// 	cipherByte, err := hex.DecodeString(cipherText)
// 	if err != nil {
// 		return out, err
// 	}
// 	sm4_out, err := sm4.Sm4Cbc([]byte(key), cipherByte, false)
// 	if err != nil {
// 		return out, err
// 	}
// 	out = string(sm4_out)
// 	return out, err
// }

func Sm3Encrypt(str string) string {
	h := sm3.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

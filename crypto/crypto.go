package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	nonceSize            = 12
	minEncryptedDataSize = nonceSize + 16 // nonce大小 + 最小加密数据大小
)

func DecryptPassword(masterKey, encryptedPass []byte) (string, error) {
	// 如果是空密码
	if len(encryptedPass) == 0 {
		return "", nil
	}

	// 尝试使用DPAPI直接解密
	decrypted, err := DecryptWithDPAPI(encryptedPass) // 修改这里，使用大写的函数名
	if err == nil {
		return string(decrypted), nil
	}

	// 如果DPAPI解密失败，使用AES-GCM解密
	if len(encryptedPass) < minEncryptedDataSize {
		return "", errors.New("encrypted password too short")
	}

	nonce := encryptedPass[3 : 3+nonceSize]
	payload := encryptedPass[3+nonceSize:]

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

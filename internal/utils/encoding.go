package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/n-creativesystem/rbns/utilsconv"
	"golang.org/x/crypto/pbkdf2"
)

func RandomString(n int, alphabets ...byte) (string, error) {
	if v, err := RandomBytes(n, alphabets...); err != nil {
		return "", err
	} else {
		return utilsconv.BytesToString(v), nil
	}
}

func RandomBytes(n int, alphabets ...byte) ([]byte, error) {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes, nil
}

func EncodePassword(password string, salt string) (string, error) {
	newPasswd := pbkdf2.Key([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(newPasswd), nil
}

func DecodeBasicAuthHeader(header string) (string, string, error) {
	var code string
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		code = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", "", err
	}

	userAndPass := strings.SplitN(string(decoded), ":", 2)
	if len(userAndPass) != 2 {
		return "", "", errors.New("invalid basic auth header")
	}

	return userAndPass[0], userAndPass[1], nil
}

type GCM struct {
	Key        []byte
	CipherText string
}

func EncryptByGCM(plainText string) (GCM, error) {
	var result GCM
	key, err := RandomBytes(32)
	if err != nil {
		return result, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return result, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return result, err
	}
	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
	cipherText = append(nonce, cipherText...)
	result.Key = key
	result.CipherText = utilsconv.BytesToBase64(cipherText)
	return result, nil
}

func DecryptGCN(encryptGCM GCM) (string, error) {
	cipherByte, err := utilsconv.Base64ToByte(encryptGCM.CipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encryptGCM.Key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := cipherByte[:gcm.NonceSize()]
	plainByte, err := gcm.Open(nil, nonce, cipherByte[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return utilsconv.BytesToString(plainByte), nil
}

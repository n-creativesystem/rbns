package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCM(t *testing.T) {
	gcm, err := EncryptByGCM("12345")
	if !assert.NoError(t, err) {
		return
	}
	decryptText, err := DecryptGCN(gcm)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "12345", decryptText)
}

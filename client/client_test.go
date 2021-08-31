package client_test

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/client"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client, err := client.New("localhost:8888")
	if !assert.NoError(t, err) {
		return
	}
	defer client.Close()
	assert.Equal(t, "localhost:8888", client.Target())
	ctx := context.Background()
	if resp, _ := client.Permissions(ctx).FindAll(nil); resp != nil {
		assert.True(t, len(resp.Permissions) > 0)
	}
}

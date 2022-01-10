package restserver

import (
	"encoding/json"
	"testing"

	"github.com/jmespath/go-jmespath"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestJMESPath(t *testing.T) {
	var jsondata = []byte(`{"attr": {"http://schemas.auth0.com/https://ncs-kubernetes;com;groups": ["admin"] }}`)
	var data interface{}
	_ = json.Unmarshal(jsondata, &data)
	exp := "contains(attr.\"http://schemas.auth0.com/https://ncs-kubernetes;com;groups\", 'admin') && 'Admin' || ''"
	result, err := jmespath.Search(exp, data)
	assert.NoError(t, err)
	logrus.Info(result)
}

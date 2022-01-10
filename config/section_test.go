package config

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSection(t *testing.T) {
	const raw = `
test:
  aaa:
    key1: value1
  bbb:
    key2: value2
`
	v := viper.New()
	v.SetConfigType("yaml")
	_ = v.ReadConfig(bytes.NewBufferString(raw))
	section := &Section{key: ""}
	sec := section.Section("test.aaa")
	assert.Equal(t, "value1", sec.Key("key1").String())
	sec = section.Section("test.bbb")
	assert.Equal(t, "value2", sec.Key("key2").String())
}

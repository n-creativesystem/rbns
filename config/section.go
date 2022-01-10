package config

import (
	"bytes"
	"os"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Section struct {
	key string
}

func (s *Section) Section(key string) *Section {
	buf := strings.Builder{}
	if s.key != "" {
		buf.WriteString(s.key)
	}
	if key != "" {
		buf.WriteRune('.')
		buf.WriteString(key)
	}
	return &Section{
		key: buf.String(),
	}
}

func (s *Section) Key(key string) Key {
	buf := strings.Builder{}
	if s.key != "" {
		buf.WriteString(s.key)
	}
	if key != "" {
		buf.WriteRune('.')
		buf.WriteString(key)
	}
	return Key{
		key: buf.String(),
	}
}

type Key struct {
	key string
}

func (k Key) Int() int {
	return viper.GetInt(k.key)
}

func (k Key) String() string {
	return viper.GetString(k.key)
}

func (k Key) Bool() bool {
	return viper.GetBool(k.key)
}

func (k Key) StringExpand() string {
	return os.ExpandEnv(k.String())
}

func (k Key) MustBool(defaultVal ...bool) bool {
	val, err := cast.ToBoolE(viper.Get(k.key))
	if len(defaultVal) > 0 && err != nil {
		viper.Set(k.key, defaultVal[0])
		return defaultVal[0]
	}
	return val
}

func (k Key) MustString(defaultValue string) string {
	value := k.String()
	if len(value) == 0 {
		viper.Set(k.key, defaultValue)
		return defaultValue
	}
	return value
}

func (k Key) Strings(delim string) []string {
	str := k.String()
	if len(str) == 0 {
		return []string{}
	}

	runes := []rune(str)
	vals := make([]string, 0, 2)
	var buf bytes.Buffer
	escape := false
	idx := 0
	for {
		if escape {
			escape = false
			if runes[idx] != '\\' && !strings.HasPrefix(string(runes[idx:]), delim) {
				buf.WriteRune('\\')
			}
			buf.WriteRune(runes[idx])
		} else {
			if runes[idx] == '\\' {
				escape = true
			} else if strings.HasPrefix(string(runes[idx:]), delim) {
				idx += len(delim) - 1
				vals = append(vals, strings.TrimSpace(buf.String()))
				buf.Reset()
			} else {
				buf.WriteRune(runes[idx])
			}
		}
		idx++
		if idx == len(runes) {
			break
		}
	}

	if buf.Len() > 0 {
		vals = append(vals, strings.TrimSpace(buf.String()))
	}

	return vals
}

func (k Key) StringSlice() []string {
	return viper.GetStringSlice(k.key)
}

func (k Key) StringMap() map[string]interface{} {
	return viper.GetStringMap(k.key)
}

package mysql

import (
	"fmt"
	"strings"
)

func concatDSN(settings []Setting) string {
	s := ""
	for _, f := range settings {
		s = f(s)
	}
	return strings.TrimRight(s, "&")
}

type Setting func(string) string

func settingString(source string, param string, value interface{}) string {
	if "" == value {
		return ""
	}
	return fmt.Sprintf(cDSNFormat, source, param, value)
}

func SetCharset(value string) Setting {
	return func(source string) string {
		return settingString(source, "charset", value)
	}
}

func SetLocal(local string) Setting {
	return func(source string) string {
		return settingString(source, "loc", local)
	}
}

func SetQueryValue(key string, value string) Setting {
	return func(source string) string {
		return settingString(source, key, value)
	}
}

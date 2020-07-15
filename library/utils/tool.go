package utils

import (
	"filesrv/conf"
	"strings"
)

func IsImage(ex string) bool {
	exName := strings.ToUpper(ex)
	for _, v := range conf.ImageExName {
		if v == exName {
			return true
		}
	}
	return false
}

package utils

import (
	"net/url"
)

func SanatizeKey(key string) string {
	sanatizedKey := url.PathEscape(key)
	return sanatizedKey
}


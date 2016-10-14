package lib

import "net/url"


func DecodeUrlComponent(u string) (string, error) {
	return url.QueryUnescape(u)
}


func EncodeUrlComponent(u string) string {
	return url.QueryEscape(u)
}
package http

import (
	"net/http"
	"strings"
)

func GetBearToken(header http.Header) string {
	//get token from header
	token := header.Get("Authorization")
	if len(token) == 0 {
		//token not exist
		return ""
	}
	// check token
	fields := strings.Fields(token)
	if len(fields) < 2 {
		//token error, no auth
		return ""
	}
	return fields[1]
}

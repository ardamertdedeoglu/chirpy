package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth_header := headers.Get("Authorization")
	if auth_header == "" {
		return "", errors.New("Authorization header not found")
	}

	return strings.TrimPrefix(auth_header, "ApiKey "), nil
}

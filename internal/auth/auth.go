package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetToken(header http.Header, keyType string) (string, error) {
	result := strings.Split(header.Get("Authorization"), " ")

	if len(result) != 2 {
		return "", fmt.Errorf("wrong authorization header format")
	}

	if strings.Trim(result[0], " ") != keyType {
		return "", fmt.Errorf("missing " + keyType)
	}

	return strings.Trim(result[1], " "), nil
}

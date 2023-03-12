package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetUrlParam(r *http.Request, param string) (string, error) {
	value := mux.Vars(r)[param]
	if value == "" {
		return "", fmt.Errorf("parameter missing in url")
	}
	return value, nil
}

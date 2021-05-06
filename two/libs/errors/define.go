package _err

import "net/http"

var (
	NotFound = New(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	Internal = New(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
)
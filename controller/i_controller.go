package controller

import "net/http"

// IController interface controller
type IController interface {
	SendHTTPResponse(w http.ResponseWriter, model interface{})
}

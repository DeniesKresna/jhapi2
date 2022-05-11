package controller

import (
	"encoding/json"
	"net/http"

	"github.com/DeniesKresna/jhapi2/service/user"

	iutil "github.com/DeniesKresna/jhapi2/utils"

	"github.com/go-playground/validator/v10"
)

// Controller struct
type Controller struct {
	validator *validator.Validate
	userSvc   user.IService
	util      iutil.IUtils
}

// Provide is contructor
func Provide(v *validator.Validate, userSvc user.IService, util iutil.IUtils) IController {
	return &Controller{
		validator: v,
		userSvc:   userSvc,
		util:      util,
	}
}

// SendHTTPResponse common func
func (u *Controller) SendHTTPResponse(w http.ResponseWriter, model interface{}) {

	// Set CORS headers for the main request.
	w.Header().Add("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(model)
}

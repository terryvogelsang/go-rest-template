package router

import (
	"encoding/json"
	"mycnc-rest-api/models"
	"mycnc-rest-api/utils"
	"net/http"

	"github.com/gorilla/mux"
	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

// CreateRegatta : Create regatta and stores it in DB
func CreateRegatta(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Parse Request Body
	var regattaCreateRequestBody models.RegattaCreateRequestBody
	err := json.NewDecoder(r.Body).Decode(&regattaCreateRequestBody)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	// Add user to DB
	regatta, err := env.GORM.CreateRegatta(&regattaCreateRequestBody)

	if err != nil {

		// User already exists error
		if match := utils.CaseInsensitiveContains(err.Error(), "duplicate entry"); match {
			return customhttpresponse.CodeAlreadyExists, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(regatta, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// ReadRegatta : Read regatta from DB
func ReadRegatta(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Parse regattaID from URL
	params := mux.Vars(r)
	regattaID := params["regattaID"]

	// Add user to DB
	regatta, err := env.GORM.ReadRegattaFromID(regattaID)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(regatta, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

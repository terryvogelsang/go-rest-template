package router

import (
	"encoding/json"
	"errors"
	"mycnc-rest-api/auth"
	"mycnc-rest-api/models"
	middlewares "mycnc-rest-api/router/middlewares"
	"mycnc-rest-api/utils"
	"net/http"

	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

// CreateUser : Create user and stores it in DB
func CreateUser(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Parse Request Body
	var userCreateRequestBody models.UserCreateRequestBody
	err := json.NewDecoder(r.Body).Decode(&userCreateRequestBody)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	hashedPassword, err := auth.HashPassword(userCreateRequestBody.Password)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	userCreateRequestBody.Password = hashedPassword

	// Add user to DB
	user, err := env.GORM.CreateUser(&userCreateRequestBody)

	if err != nil {

		// User already exists error
		if match := utils.CaseInsensitiveContains(err.Error(), "duplicate entry"); match {
			return customhttpresponse.CodeAlreadyExists, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(user.ID, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// ReadUser : Read user informations from DB
func ReadUser(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	userID := r.Context().Value(middlewares.ContextUserKey).(string)

	user, err := env.GORM.ReadUserFromID(userID)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(user, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// ReadUserBoatsInfos : Read user boats informations from DB
func ReadUserBoatsInfos(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get existing user from DB
	userID := r.Context().Value(middlewares.ContextUserKey).(string)
	user, err := env.GORM.ReadUserFromID(userID)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	// Get user boats from DB
	boats, err := env.GORM.ReadUserBoatsInfos(user)

	if err != nil {

		if env.GORM.IsRecordNotFoundError(err) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(boats, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// UpdateUser : Update user and stores it in DB
func UpdateUser(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get existing user from DB
	userID := r.Context().Value(middlewares.ContextUserKey).(string)
	user, err := env.GORM.ReadUserFromID(userID)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	// Parse Request Body
	var userUpdateRequest models.UserUpdateRequestBody
	err = json.NewDecoder(r.Body).Decode(&userUpdateRequest)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	// Update user in DB
	err = env.GORM.UpdateUserInfos(user, &userUpdateRequest)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(user, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// UpdateUserPassword : Update user password and stores it in DB
func UpdateUserPassword(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	userID := r.Context().Value(middlewares.ContextUserKey).(string)

	// Parse Request Body
	var changePasswordRequest = &models.ChangePasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(&changePasswordRequest)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	// Check current password
	user, err := env.GORM.ReadUserFromID(userID)

	if err != nil {

		if env.GORM.IsRecordNotFoundError(err) {
			return customhttpresponse.CodeBadLogin, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	// user.Password represents the hashed pasword from DB
	if auth.CheckPasswordHash(changePasswordRequest.OldPassword, user.Password) {

		// Hash new password
		newHashedPassword, err := auth.HashPassword(changePasswordRequest.NewPassword)

		if err != nil {
			return customhttpresponse.CodeInternalError, err
		}

		// Update password in DB
		err = env.GORM.UpdateUserPassword(newHashedPassword)

		if err != nil {
			return customhttpresponse.CodeInternalError, err
		}

		responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
		customhttpresponse.WriteResponse(nil, responseDetails, w)

		return customhttpresponse.CodeSuccess, nil
	}

	return customhttpresponse.CodeBadLogin, errors.New("Wrong old password")
}

// UpdateUserBoats : Add boat to user in DB
func UpdateUserBoats(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get existing user from DB
	userID := r.Context().Value(middlewares.ContextUserKey).(string)
	user, err := env.GORM.ReadUserFromID(userID)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	// Parse Request Body and create boat struct
	var boatCreateRequestBody models.BoatCreateRequestBody
	err = json.NewDecoder(r.Body).Decode(&boatCreateRequestBody)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	err = env.GORM.UpdateUserBoats(user, &boatCreateRequestBody)

	if err != nil {
		if env.GORM.IsRecordNotFoundError((err)) {
			return customhttpresponse.CodeDoesNotExist, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(user.Boats, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

// DeleteUser : Delete user from DB
func DeleteUser(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	var user = &models.User{}

	userID := r.Context().Value(middlewares.ContextUserKey).(string)

	// Prevent ID Overwriting from request
	user.ID = userID

	// Delete user from DB
	err := env.GORM.DeleteUser(user)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)
	customhttpresponse.WriteResponse(nil, responseDetails, w)

	return customhttpresponse.CodeSuccess, nil
}

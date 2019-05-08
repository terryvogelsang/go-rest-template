package router

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mycnc-rest-api/auth"
	"mycnc-rest-api/models"
	middlewares "mycnc-rest-api/router/middlewares"
	"mycnc-rest-api/utils"
	"net/http"
	"strings"
	"time"

	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

// CreateSession : Log user in and generate session
func CreateSession(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Parse Request Body
	var credentials models.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		return customhttpresponse.CodeInvalidJSON, err
	}

	// Check match in DB
	user, err := env.GORM.ReadUserFromEmail(credentials.Email)

	if err != nil {

		if env.GORM.IsRecordNotFoundError(err) {
			return customhttpresponse.CodeBadLogin, err
		}

		return customhttpresponse.CodeInternalError, err
	}

	// user.Password represents the hashed pasword from DB
	if auth.CheckPasswordHash(credentials.Password, user.Password) {

		userStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisUserStoragePrefix, user.ID, models.RedisUserStorageSessionSuffix)
		existingSession, err := env.Redis.Get(userStorageKey)
		existingSessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, existingSession, models.RedisSessionStorageUserIDSuffix)

		if err != nil && !utils.CaseInsensitiveContains(err.Error(), "error getting key user") {

			return customhttpresponse.CodeInternalError, err
		}

		// If a session already exists for this user, revoke it
		if string(existingSession) != "" {
			env.Redis.Delete(existingSessionStorageKey)
		}

		// Generate Session Token
		randomBytes, err := utils.GenerateCryptoRandomBytes(16)

		if err != nil {
			return customhttpresponse.CodeInternalError, err
		}

		sessionToken := strings.ToUpper(hex.EncodeToString(randomBytes))
		sessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, sessionToken, models.RedisSessionStorageUserIDSuffix)

		// Init Redis transaction to add datas to User and Session storage
		transactionCommands := []models.RedisCommand{
			models.RedisCommand{
				Command: "SETEX",
				Args:    []interface{}{sessionStorageKey, models.TokenExpirationInMinutes * 60, []byte(user.ID)},
			},
			models.RedisCommand{
				Command: "SET",
				Args:    []interface{}{userStorageKey, []byte(sessionToken)},
			},
		}

		_, err = env.Redis.Multi(transactionCommands)

		if err != nil {
			return customhttpresponse.CodeInternalError, err
		}

		// Cookie expiration fixed to 30 minutes
		maxCookieAge := (time.Duration(models.TokenExpirationInMinutes) * time.Minute) / time.Second

		tokenCookie := http.Cookie{

			// __Host cookie prefix signals to the browser that both the Path=/
			// and Secure attributes are required, and at the same time
			// + that the Domain attribute must not be present.
			// see https://resources.infosecinstitute.com/cookies-httponly-flag-problem-browsers
			// TO DO : Change this to __Host-session
			Name: "session",

			// Path set to root
			Path: "/",

			// TODO : Uncomment this
			// Only transmit on encrypted connection
			// Secure: true,

			// Cannot be accessed by JavaScript
			HttpOnly: true,

			// Prevents cookie from being attached to cross-origin requests
			SameSite: http.SameSiteStrictMode,

			// Actual value of cookie
			Value: sessionToken,

			// Set Expiration cookie client side (IE Browser will treat this cookie as a session cookie [Expires when closing tab])
			MaxAge: int(maxCookieAge),
		}

		// Set cookie to response
		http.SetCookie(w, &tokenCookie)

		// Return response
		responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)

		customhttpresponse.WriteResponse(
			struct {
				Session string `json:"session"`
			}{
				sessionToken,
			}, responseDetails, w,
		)

		return customhttpresponse.CodeSuccess, nil
	}

	return customhttpresponse.CodeBadLogin, errors.New("Invalid credentials")
}

// UpdateSession : Refresh session
func UpdateSession(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	userID := r.Context().Value(middlewares.ContextUserKey).(string)

	// Get token from cookies
	c, err := r.Cookie("session")

	if err != nil && c.Value != "" {
		return "", err
	}

	existingSessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, c.Value, models.RedisSessionStorageUserIDSuffix)
	userStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisUserStoragePrefix, userID, models.RedisUserStorageSessionSuffix)

	// Delete key pair from session storage
	err = env.Redis.Delete(existingSessionStorageKey)

	if err != nil {
		return customhttpresponse.CodeInternalError, errors.New("Could not delete session in Redis")
	}

	// Generate new session token
	randomBytes, err := utils.GenerateCryptoRandomBytes(16)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	sessionToken := strings.ToUpper(hex.EncodeToString(randomBytes))
	sessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, sessionToken, models.RedisSessionStorageUserIDSuffix)

	// Init Redis transaction to add datas to User and Session storage
	transactionCommands := []models.RedisCommand{
		models.RedisCommand{
			Command: "SETEX",
			Args:    []interface{}{sessionStorageKey, models.TokenExpirationInMinutes * 60, []byte(userID)},
		},
		models.RedisCommand{
			Command: "SET",
			Args:    []interface{}{userStorageKey, []byte(sessionToken)},
		},
	}

	_, err = env.Redis.Multi(transactionCommands)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	// Cookie expiration fixed to 30 minutes
	maxCookieAge := (time.Duration(models.TokenExpirationInMinutes) * time.Minute) / time.Second

	tokenCookie := http.Cookie{

		// __Host cookie prefix signals to the browser that both the Path=/
		// and Secure attributes are required, and at the same time
		// + that the Domain attribute must not be present.
		// see https://resources.infosecinstitute.com/cookies-httponly-flag-problem-browsers
		// TO DO : Change this to __Host-session
		Name: "session",

		// Path set to root
		Path: "/",

		// TODO : Uncomment this
		// Only transmit on encrypted connection
		// Secure: true,

		// Cannot be accessed by JavaScript
		HttpOnly: true,

		// Prevents cookie from being attached to cross-origin requests
		SameSite: http.SameSiteStrictMode,

		// Actual value of cookie
		Value: sessionToken,

		// Set Expiration cookie client side (IE Browser will treat this cookie as a session cookie [Expires when closing tab])
		MaxAge: int(maxCookieAge),
	}

	// Set cookie to response
	http.SetCookie(w, &tokenCookie)

	// Return response
	responseDetails := customhttpresponse.NewResponseDetails(env.Config.Service, utils.GetCurrentFuncName(), customhttpresponse.CodeSuccess)

	customhttpresponse.WriteResponse(
		struct {
			Session string `json:"session"`
		}{
			sessionToken,
		}, responseDetails, w,
	)

	return customhttpresponse.CodeSuccess, nil
}

// DeleteSession : Log user out and delete session from storage
func DeleteSession(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	userID := r.Context().Value(middlewares.ContextUserKey).(string)

	// Get token from cookies
	c, err := r.Cookie("session")

	if err != nil && c.Value != "" {
		return "", err
	}

	existingSessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, c.Value, models.RedisSessionStorageUserIDSuffix)
	userStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisUserStoragePrefix, userID, models.RedisUserStorageSessionSuffix)

	// Delete key pair from session storage
	err = env.Redis.Delete(existingSessionStorageKey)

	if err != nil {
		return customhttpresponse.CodeInternalError, errors.New("Could not delete session in Redis")
	}

	// Set session as nil in user storage
	err = env.Redis.Set(userStorageKey, nil)

	if err != nil {
		return customhttpresponse.CodeInternalError, err
	}

	return customhttpresponse.CodeSuccess, nil
}

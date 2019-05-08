package router

import (
	errors "errors"
	"fmt"
	models "mycnc-rest-api/models"
	http "net/http"

	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

type (
	ContextKey string
)

const (
	ContextUserKey ContextKey = "userID"
)

var (
	serviceVersion   = "/v1"
	userRoute        = serviceVersion + "/user"
	authSessionRoute = serviceVersion + "/auth/session"

	// These routes are publicly accessible without authentication
	unauthenticatedRoutes = map[string]map[string]bool{

		// POST /v1/user (Create User)
		userRoute: map[string]bool{
			http.MethodPost: true,
		},

		// POST /v1/auth/session (Create session - Login)
		authSessionRoute: map[string]bool{
			http.MethodPost: true, // Create session (LOGIN)
		},
	}
)

// AuthMiddleware : Check session token
func AuthMiddleware(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get request path and request method
	reqPath := r.URL.Path
	reqMethod := r.Method

	// If route is whitelisted, check method. If it matches, don't do the authentication check and immediately forward the request
	if unauthenticatedRoutes[reqPath] != nil {
		if unauthenticatedRoutes[reqPath][reqMethod] {
			return "", nil
		}
	}

	// Get token from cookies
	c, err := r.Cookie("session")

	// If no token, but authentication is needed, don't forward the request
	if err != nil {
		return "", err
	}

	t := c.Value

	if t == "" {
		return "", errors.New("Empty session string")
	}

	sessionStorageKey := fmt.Sprintf("%s:%s:%s", models.RedisSessionStoragePrefix, t, models.RedisSessionStorageUserIDSuffix)

	// Get associated UserID
	userID, err := env.Redis.Get(sessionStorageKey)

	if err != nil {
		return "", err
	}

	return string(userID), nil
}

// SessionExistsInStorage : Check if session exists in Redis
func SessionExistsInStorage(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get token from cookies
	c, err := r.Cookie("session")

	// If no token, but authentication is needed, don't forward the request
	if err != nil {
		return "", err
	}

	// Check if key exists in Redis
	_, err = env.Redis.Exists(c.Value)

	if err != nil {
		return customhttpresponse.CodeDoesNotExist, err
	}

	return "", nil
}

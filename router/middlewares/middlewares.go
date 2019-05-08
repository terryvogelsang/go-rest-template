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
	serviceVersion    = "/v1"
	userRoute         = serviceVersion + "/user"
	userPasswordRoute = userRoute + "/password"
	userBoatsRoute    = userRoute + "/boats"

	regattaRoute = serviceVersion + "/regatta"

	authSessionRoute = serviceVersion + "/auth/session"

	// Whitelisted routes are publicly accessibe without authentication
	whitelistedRoutes = map[string]map[string]bool{
		userRoute: map[string]bool{
			http.MethodPost: true, // Create User
		},
		authSessionRoute: map[string]bool{
			http.MethodPost: true, // Create session (LOGIN)
		},
	}

	// targets = map[string]*ProxyRequest{

	// 	userRoute: &ProxyRequest{RequestPath: userRoute, IsMethodProtected: map[string]bool{
	// 		http.MethodPost:   false, // Create user (Unprotected)
	// 		http.MethodGet:    true,  // Read user infos (Protected)
	// 		http.MethodPut:    true,  // Update user infos (Protected)
	// 		http.MethodDelete: true,  // Delete user infos (Protected)
	// 	}},

	// 	userPasswordRoute: &ProxyRequest{RequestPath: userPasswordRoute, IsMethodProtected: map[string]bool{
	// 		http.MethodPut: true, // Update user password (Protected)
	// 	}},

	// 	userBoatsRoute: &ProxyRequest{RequestPath: userBoatsRoute, IsMethodProtected: map[string]bool{
	// 		http.MethodPut: true, // Update user boats (Protected)
	// 		http.MethodGet: true, // Get user boats infos (Protected)
	// 	}},

	// 	authSessionRoute: &ProxyRequest{RequestPath: authSessionRoute, IsMethodProtected: map[string]bool{
	// 		http.MethodPost:   false, // Create session (Unprotected)
	// 		http.MethodPut:    true,  // Refresh session (Protected)
	// 		http.MethodDelete: true,  // Delete session (Protected)
	// 	}},
	// 	regattaRoute: &ProxyRequest{RequestPath: regattaRoute, IsMethodProtected: map[string]bool{
	// 		http.MethodPost: true, // Create regatta (Protected)
	// 		http.MethodGet:  true, // Read regatta (Protected)

	// 	}},
	// }
)

// ProxyRequest : Received Request from Client
type ProxyRequest struct {
	RequestPath       string
	IsMethodProtected map[string]bool
}

// AuthMiddleware : Check session token
func AuthMiddleware(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error) {

	// Get request path and request method
	reqPath := r.URL.Path
	reqMethod := r.Method

	// If route is whitelisted, check method. If it matches, don't do the authentication check and immediately forward the request
	if whitelistedRoutes[reqPath] != nil {
		fmt.Println("yoli")
		if whitelistedRoutes[reqPath][reqMethod] {
			fmt.Println("yolo")
			return "", nil
		}
	}

	// If method for path is whitelisted, do not check auth

	// strippedPath := r.URL.Path
	// regex := regexp.MustCompile(".*[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

	// // If UUID is present remove trailing id from url
	// if regex.Match([]byte(strippedPath)) {
	// 	parts := strings.Split(r.URL.Path, "/")
	// 	strippedPath = strings.Replace(r.URL.Path, fmt.Sprintf("/%s", parts[len(parts)-1]), "", 1)
	// }

	// // If target is not registered,
	// if targets[strippedPath] == nil {
	// 	return "", errors.New("No targets configured in middleware")
	// }

	// // If the Target doesn't need authentication, don't do the authentication check and immediately forward the response
	// if !targets[strippedPath].IsMethodProtected[r.Method] {
	// 	return "", nil
	// }

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

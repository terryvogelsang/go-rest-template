package router

import (

	// Native Go Libs
	"context"
	models "mycnc-rest-api/models"
	middlewares "mycnc-rest-api/router/middlewares"
	http "net/http"
	"reflect"
	"runtime"
	"strings"

	// 3rd Party Libs
	customhttpresponse "github.com/terryvogelsang/go-custom-http-response"
)

type (
	// Handler : Custom type to work with CustomHandle wrapper
	Handler func(env *models.Env, w http.ResponseWriter, r *http.Request) (string, error)
)

type Greeter struct {
	Message string
}

// CustomHandle : Custom Handlers Wrapper for API
func CustomHandle(env *models.Env, handlers ...Handler) http.Handler {

	statusCode := ""

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		responseDetails := &customhttpresponse.ResponseDetails{}

		// Retrieve AuthMiddleware method name for response details
		action := strings.Split(runtime.FuncForPC(reflect.ValueOf(middlewares.AuthMiddleware).Pointer()).Name(), ".")[1]

		// Get UserID through authentication middleware
		userID, err := middlewares.AuthMiddleware(env, w, r)

		if err != nil {
			responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, customhttpresponse.CodeInvalidToken)
			customhttpresponse.WriteResponse(nil, responseDetails, w)
			return
		}

		// Pass UserID to request context
		ctx := context.WithValue(r.Context(), middlewares.ContextUserKey, userID)

		if err != nil {

			if statusCode == customhttpresponse.CodeValidationFailed {
				responseDetails = customhttpresponse.NewResponseDetailsWithFields(strings.Split(err.Error(), "|"), env.Config.Service, runtime.FuncForPC(reflect.ValueOf(middlewares.AuthMiddleware).Pointer()).Name(), statusCode)
			} else {
				// FIXME: Remove Debugging Mode before production
				responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, statusCode)
			}

			customhttpresponse.WriteResponse(nil, responseDetails, w)

			// We can then log error somewhere here

			return
		}

		// If auth check is successful, trigger handlers
		for _, h := range handlers {

			// Retrieve handler method name for response details
			action = strings.Split(runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), ".")[1]
			statusCode, err = h(env, w, r.WithContext(ctx))

			if err != nil {

				if statusCode == customhttpresponse.CodeValidationFailed {
					responseDetails = customhttpresponse.NewResponseDetailsWithFields(strings.Split(err.Error(), "|"), env.Config.Service, runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(), statusCode)
				} else {
					// FIXME: Remove Debugging Mode before production
					responseDetails = customhttpresponse.NewResponseDetailsWithDebug(err.Error(), env.Config.Service, action, statusCode)
				}

				customhttpresponse.WriteResponse(nil, responseDetails, w)

				// We can then log error somewhere here

				return
			}
		}

		// We can then log success somewhere here
	})
}

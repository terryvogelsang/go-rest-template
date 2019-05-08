package router

import (

	// Native Go Libs
	fmt "fmt"
	http "net/http"

	// Project Libs
	handlers "mycnc-rest-api/router/handlers"

	models "mycnc-rest-api/models"

	// 3rd Party Libs
	middlewares "mycnc-rest-api/router/middlewares"

	mux "github.com/gorilla/mux"
	cors "github.com/rs/cors"
)

// Listen : Defines all router routing rules and handlers.
// Serves the API at defined port constant.
func Listen(env *models.Env) {

	r := mux.NewRouter().StrictSlash(false)

	v1 := r.PathPrefix("/v1").Subrouter()

	// API Endpoints

	// User
	userV1 := v1.PathPrefix("/user").Subrouter()
	userV1.Handle("", handlers.CustomHandle(env, handlers.CreateUser)).Methods("POST")
	userV1.Handle("", handlers.CustomHandle(env, handlers.ReadUser)).Methods("GET")
	userV1.Handle("", handlers.CustomHandle(env, handlers.UpdateUser)).Methods("PUT")
	userV1.Handle("", handlers.CustomHandle(env, handlers.DeleteUser)).Methods("DELETE")
	userV1.Handle("/password", handlers.CustomHandle(env, handlers.UpdateUserPassword)).Methods("PUT")
	userV1.Handle("/boats", handlers.CustomHandle(env, handlers.UpdateUserBoats)).Methods("PUT")
	userV1.Handle("/boats", handlers.CustomHandle(env, handlers.ReadUserBoatsInfos)).Methods("GET")

	// Auth
	authV1 := v1.PathPrefix("/auth").Subrouter()
	authSessionV1 := authV1.PathPrefix("/session").Subrouter()
	authSessionV1.Handle("", handlers.CustomHandle(env, handlers.CreateSession)).Methods("POST")
	authSessionV1.Handle("", handlers.CustomHandle(env, handlers.UpdateSession)).Methods("PUT")
	authSessionV1.Handle("", handlers.CustomHandle(env, middlewares.SessionExistsInStorage, handlers.DeleteSession)).Methods("DELETE")

	// Regatta
	regattaV1 := v1.PathPrefix("/regatta").Subrouter()
	regattaV1.Handle("", handlers.CustomHandle(env, handlers.CreateRegatta)).Methods("POST")
	regattaV1.Handle("/{regattaID}", handlers.CustomHandle(env, handlers.ReadRegatta)).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With"},
		AllowedOrigins:   []string{"http://frontend.localhost"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	fmt.Println("Listening on port :" + fmt.Sprintf("%d", env.Config.ListeningPort))
	http.ListenAndServe(":"+fmt.Sprintf("%d", env.Config.ListeningPort), corsHandler.Handler(r))
}

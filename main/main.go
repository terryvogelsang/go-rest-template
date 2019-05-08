package main

import (

	// Native Go Libs
	fmt "fmt"
	log "log"
	os "os"

	// Project Libs
	models "mycnc-rest-api/models"
	router "mycnc-rest-api/router"
)

var (
	GlobalConfig      *models.Config
	GORMDBUser        = "root"
	GORMDBPassword    = "example"
	GORMDBName        = "mycnc"
	GORMConnectionURL = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", GORMDBUser, GORMDBPassword, GORMDBName)
	RedisHost         = "localhost"
	RedisPort         = 6379
	RedisPassword     = "example"
	RedisURL          = fmt.Sprintf("redis://%s:%d", RedisHost, RedisPort)
)

func main() {

	if os.Getenv(models.ConfigFilePathName) == "" {
		log.Fatalf(fmt.Sprintf("%s Environment variable must be set !", models.ConfigFilePathName))
	}

	gorm := models.NewGORM(GORMConnectionURL)
	redis := models.NewRedis(RedisURL, RedisPassword)

	// Add interfaces & blank config to the environment
	env := &models.Env{
		// Load here your databases communication instances
		GORM:   gorm,
		Redis:  redis,
		Config: models.Config{},
	}

	// Dynamically load config
	err := env.RefreshConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	router.Listen(env)

	defer func() {
		gorm.CloseConnection()
		redis.CloseConnection()
	}()
}

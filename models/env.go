package models

import (
	json "encoding/json"
	ioutil "io/ioutil"
	os "os"
)

var (
	GlobalConfig       *Config
	ConfigFilePathName = "MYCNC_REST_API_CONFIG_FILE_PATH"
	configFilePath     = os.Getenv(ConfigFilePathName)
)

// Env : Execution environment containing Datastore communication interfaces & Config
type Env struct {
	// Add Databases communication interfaces here
	GORM   GORMInterface
	Redis  RedisInterface
	Config Config
}

// Config : Global Config
type Config struct {
	// Add config structures here
	Service       string   `json:"service"`
	ListeningPort int      `json:"listeningPort"`
	AdminUsers    []string `json:"adminUsers"`
}

// RefreshConfig : Load current environment values in config
func (env *Env) RefreshConfig() error {

	data, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &env.Config)

	// GlobalConfig used for access in models
	GlobalConfig = &env.Config

	if err != nil {
		return err
	}

	return nil
}

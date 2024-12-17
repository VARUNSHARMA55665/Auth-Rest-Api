package resources

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

var configInstance *Config // package private singleton instance of the configuration
var singleton sync.Once    // package private singleton helper utility

func GetConfig() *viper.Viper {
	// create an instance if not available
	singleton.Do(func() {
		configInstance = &Config{viper.New()}
	})

	return configInstance.viper
}

func Start() {
	// Find and read the config file
	err := GetConfig().ReadInConfig()
	log.Printf("Start config err =%v\n", err)
	if err != nil {
		// Handle errors reading the config file
		log.Printf("Start Error Reading config = %v\n", err)
	}
}

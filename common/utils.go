package common

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// ViperEnvVariable is func to get .env file
func ViperEnvVariable(key string) string {
	//switch for reducing the number of open files (.env)
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		fmt.Println(key)
		log.Fatalf("Invalid type assertion")
	}

	return value
}

package common

import (
	"fmt"
	"log"

	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spf13/viper"
)

// ViperEnvVariable is func to get .env file
func GetEnv(key string) string {
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

func SetHash(str string) uint64 {
	h1 := fnv1a.HashString64(str)
	fmt.Println("FNV-1a hash of ", str, ":", h1)

	// Incrementally compute a hash value from a sequence of strings.
	return h1
}

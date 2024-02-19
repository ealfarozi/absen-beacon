package common

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spf13/viper"
)

var LOCAL_NAME string
var HASHED string
var REFRESH_INTERVAL int

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
	return h1
}

func GetHash() {
	HASHED = strconv.FormatUint(SetHash(GetEnv("BEACON_ID")+":"+GetUUID()), 10)
}

func GetVars() {
	LOCAL_NAME = GetEnv("LOCAL_NAME")
	rim, _ := strconv.Atoi(GetEnv("REFRESH_INTERVAL_MIN"))
	REFRESH_INTERVAL = rim
}

func GetUUID() string {
	id := uuid.New()
	return id.String()
}

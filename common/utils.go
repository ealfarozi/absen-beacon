package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spf13/viper"
	sqids "github.com/sqids/sqids-go"
)

var LOCAL_NAME string
var HASHED string
var REFRESH_INTERVAL int
var IS_STATIC string
var UUID string

type LoginRequest struct {
	NIK        int    `json:"trx_id,omitempty"`
	Code       int    `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Data       string `json:"data,omitempty"`
	SysMessage string `json:"system_message,omitempty"`
	Request    interface{}
}

type BeaconRequest struct {
	BeaconID   string `json:"beacon_id,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	ExpTimeMin int    `json:"expire_time_min,omitempty"`
	Data       string `json:"data,omitempty"`
}

type BeaconResponse struct {
	Success bool   `json:"success,omitempty"`
	Code    int    `json:"code,omitempty"`
	Data    string `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

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
	fmt.Println("Setup hash", str, h1)
	return h1
}

func GetHash() {
	HASHED = strconv.FormatUint(SetHash(UUID+":"+GetUUID()), 12)
}

func GetVars() {
	LOCAL_NAME = GetEnv("LOCAL_NAME")
	rim, _ := strconv.Atoi(GetEnv("REFRESH_INTERVAL_SEC"))
	REFRESH_INTERVAL = rim
	IS_STATIC = GetEnv("IS_STATIC")
	UUID = GetEnv("BEACON_ID")
}

func GetUUID() string {
	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{1, 2, 3}) // "86Rf07"
	return id
}

func HitAPI(url string, jsonStr []byte, method string, strToken string, timeout time.Duration) (*http.Request, *http.Response, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(jsonStr)))

	if err != nil {
		fmt.Println("error when hit URL:", url, "- err:", err.Error())
	} else {
		req.Close = true
		req.Header.Add("Content-Type", "application/json")
	}

	if strToken != "" {
		req.Header.Add("Authorization", strToken)
	}

	tr := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 500,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: timeout * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(time.Second * timeout)
		fmt.Println("error when hit URL:", url, "- err:", err.Error())

		return req, resp, nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return req, resp, body, nil
}

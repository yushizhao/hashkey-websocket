package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	SecretKey  string
	ApiKeyHMAC string
	Host       string
	Symbol     string
}

var (
	config   Config
	authType string
	path     string
	hmacStr  string
	baseURL  string
)

func Init(configPath *string) error {
	jsonFile, err := os.Open(*configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}

	authType = "HMAC"
	path = "/APITradeWS/v1/messages"
	hmacStr = SHA256HMAC([]byte("WSS/APITradeWS/v1/messages"), config.SecretKey)
	baseURL = config.Host + path
	return nil
}

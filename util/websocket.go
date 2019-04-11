package util

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var c *websocket.Conn
var err error

func InitWS() error {
	ticker := time.NewTicker(15 * time.Second)

	//c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	c, _, err = dialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
		return err
	} else {
		log.Printf("connecting to %s", u.String())
	}

	SubPrivate()
	SubDepth()
	SubTicker()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("recv error:", err)
				return
			} else {
				log.Printf("recv: %s", message)
				// onMessage(message)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// heartbeat mechanism
				// 心跳机制
				err = c.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Println("ping:", err)
					return
				}
			}
		}
	}()

	return nil
}

type RequestMessage struct {
	Type    string              `json:"type"`
	Channel map[string][]string `json:"channel"`
}

//https://hashkeypro.github.io/api-spec/#3-3-2-level2-market-data-emsp-level2
//{ "type": "subscribe", "channel": {"level2@10":["ETH-BTC"]} }
func SubDepth() {
	publicChannel := make(map[string][]string)
	publicChannel["level2@10"] = []string{config.Symbol}
	err := subscribe(publicChannel)
	if err != nil {
		log.Println(err.Error())
	}
}

// Private Message Flow
// { "type": "subscribe", "channel": {"ticker":["ETH-BTC", "CYB-BTC"]} }
func SubTicker() {
	publicChannel := make(map[string][]string)
	publicChannel["ticker"] = []string{config.Symbol}
	err := subscribe(publicChannel)
	if err != nil {
		log.Println(err.Error())
	}
}

// Private Message Flow
// { "type": "subscribe", "channel": {"user":[API-KEY, API-SIGNATURE, AUTHTYPE]} }
func SubPrivate() {
	privateChannel := make(map[string][]string)
	privateChannel["user"] = []string{config.ApiKeyHMAC, hmacStr, authType}
	err := subscribe(privateChannel)
	if err != nil {
		log.Println(err.Error())
	}
}

func subscribe(channel map[string][]string) error {
	message := RequestMessage{
		Type:    "subscribe",
		Channel: channel,
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Println("write:", err)
		return err
	}

	return nil
}

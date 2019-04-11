package main

import (
	"flag"
	"log"

	"github.com/yushizhao/hashkey-websocket/util"
)

func main() {
	var configPath = flag.String("path", "config.json", "the file")

	flag.Parse()

	err := util.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	util.InitWS()

	done := make(chan bool, 1)
	<-done
}

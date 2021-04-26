package main

import (
	"os"

	"github.com/desertbit/grumble"
	log "github.com/go-fastlog/fastlog"
	"github.com/xiaohuiluo/blockchain-sim/cmd"
)

func main() {

	logFile, err := os.OpenFile("simulate.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetOutput(logFile)

	grumble.Main(cmd.Cli)
}

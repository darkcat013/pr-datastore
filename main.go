package main

import (
	"math/rand"
	"time"

	"github.com/darkcat013/pr-datastore/protocols"
	"github.com/darkcat013/pr-datastore/utils"
)

func main() {
	rand.Seed(time.Now().UnixMilli())
	utils.InitializeLogger()
	go protocols.StartUdpServer()
	go protocols.StartUdpClient()
	go protocols.StartTcpServer()

	protocols.StartHttp()
}

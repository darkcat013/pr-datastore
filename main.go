package main

import (
	"github.com/darkcat013/pr-datastore/protocols"
	"github.com/darkcat013/pr-datastore/utils"
)

func main() {
	utils.InitializeLogger()
	go protocols.StartUdpServer()
	go protocols.StartUdpClient()
	go protocols.StartTcpServer()
	protocols.StartHttp()
}

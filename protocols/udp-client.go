package protocols

import (
	"encoding/json"
	"net"
	"time"

	"github.com/darkcat013/pr-datastore/config"
	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

func StartUdpClient() {
	udpBroadcastAddress := constants.UDP_BROADCAST_IP + constants.UDP_BROADCAST_PORT
	addr, err := net.ResolveUDPAddr("udp4", udpBroadcastAddress)
	if err != nil {
		utils.Log.Errorw("UDP Client | Error resolving address", "address", udpBroadcastAddress, "error", err.Error())
		return
	}

	go broadcastConfig(addr)
	go checkConnections()
}

func broadcastConfig(addr *net.UDPAddr) {
	for {
		time.Sleep(constants.UDP_BROADCAST_DELAY)

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error dialing", "address", addr, "error", err.Error())
			return
		}
		defer conn.Close()

		config := domain.ConfigUdp{Name: config.NAME, Port: config.TCP_PORT}
		body, err := json.Marshal(config)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error encoding config to JSON", "error", err.Error())
			return
		}

		_, err = conn.Write(body)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error broadcasting config", "error", err.Error())
			return
		}
	}
}

func checkConnections() {
	for {
		time.Sleep(constants.UDP_BROADCAST_DELAY)
		for connName, connection := range datastore.Connections {
			if time.Now().Unix()-connection.LastUpdated > constants.CONNECTION_ALIVE_TIME {
				utils.Log.Warnw("UDP Client | A connection is dead", "connection", connName)
				delete(datastore.Connections, connName)
			}
		}
	}
}

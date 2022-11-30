package protocols

import (
	"encoding/json"
	"net"

	"github.com/darkcat013/pr-datastore/config"
	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

func StartUdpServer() {
	udpBroadcastAddress := constants.UDP_BROADCAST_IP + constants.UDP_BROADCAST_PORT
	addr, err := net.ResolveUDPAddr("udp4", udpBroadcastAddress)
	if err != nil {
		utils.Log.Errorw("UDP Server | Error resolving address", "address", udpBroadcastAddress, "error", err.Error())
		return
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		utils.Log.Errorw("UDP Server | Error listening on address", "address", addr.String(), "error", err.Error())
		return
	}

	utils.Log.Infow("UDP Server | Start listening on ", "address", udpBroadcastAddress)

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			utils.Log.Errorw("UDP Server | Error reading", "error", err.Error())
			continue
		}

		buffer = buffer[:n]
		var configUdp domain.ConfigUdp
		err = json.Unmarshal(buffer, &configUdp)
		if err != nil {
			utils.Log.Errorw("UDP Server | Error decoding config JSON", "error", err.Error(), "info", string(buffer))
			continue
		}

		if configUdp.Name != config.NAME {
			utils.Log.Infow("UDP Server | Received config", "config", configUdp)

			UpdateTcpConnection(configUdp)
		}
	}
}

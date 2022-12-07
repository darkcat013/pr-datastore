package protocols

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"github.com/darkcat013/pr-datastore/config"
	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

var udpAddr *net.UDPAddr

func StartUdpClient() {
	var err error
	udpBroadcastAddress := constants.UDP_BROADCAST_IP + constants.UDP_BROADCAST_PORT
	udpAddr, err = net.ResolveUDPAddr("udp4", udpBroadcastAddress)
	if err != nil {
		utils.Log.Errorw("UDP Client | Error resolving address", "address", udpBroadcastAddress, "error", err.Error())
		return
	}

	go broadcastConfig()
	go checkConnections()
}

func broadcastConfig() {
	for {
		time.Sleep(constants.UDP_BROADCAST_DELAY)

		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error dialing", "address", udpAddr, "error", err.Error())
			return
		}
		defer conn.Close()

		allKeys := datastore.GetAllKeys()
		dataIds := make(map[string]bool)
		for i := 0; i < len(allKeys); i++ {
			dataIds[allKeys[i]] = true
		}

		cfg := domain.UdpAction{
			Action:   constants.UDP_CONFIG,
			Name:     config.NAME,
			Port:     config.TCP_PORT,
			IsLeader: datastore.IsLeader,
			DataIds:  dataIds,
		}
		body, err := json.Marshal(cfg)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error encoding config to JSON", "error", err.Error())
			return
		}

		_, err = conn.Write(body)
		if err != nil {
			utils.Log.Errorw("UDP Client | Error broadcasting config", "error", err.Error())
			return
		}
		utils.Log.Infow("UDP Client | Broadcasted config", "config", cfg)
	}
}

func checkConnections() {
	for {
		time.Sleep(constants.CONNECTION_ALIVE_TIME * time.Second)
		utils.Log.Infow("UDP Client | Start checking connections", "connections", datastore.Connections)

		leaderPresent := datastore.IsLeader
		for connName, connection := range datastore.Connections {
			if time.Now().Unix()-connection.LastUpdated > constants.CONNECTION_ALIVE_TIME {
				utils.Log.Warnw("UDP Client | A connection is dead", "connection", connName)
				if connection.IsLeader {
					datastore.HasLeader = false
					utils.Log.Warnw("UDP Client | Leader is dead", "connection", connName)
				}
				delete(datastore.Connections, connName)
			} else if connection.IsLeader {
				leaderPresent = true
			}
		}
		if !datastore.HasLeader || !leaderPresent {
			utils.Log.Warnw("UDP Client | Leader is either not present or dead, start election")
			go startElection()
		}
	}
}

func startElection() {

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		utils.Log.Errorw("UDP Client | Error dialing", "address", udpAddr, "error", err.Error())
		return
	}
	defer conn.Close()

	names := make([]string, 0, len(datastore.Connections))

	for k := range datastore.Connections {
		names = append(names, k)
	}

	var vote = ""
	if len(names) <= 0 {
		vote = config.NAME
	} else {
		vote = names[rand.Intn(len(names))]
	}

	electAction := domain.UdpAction{
		Action: constants.UDP_ELECT,
		Name:   vote,
	}

	body, err := json.Marshal(electAction)
	if err != nil {
		utils.Log.Errorw("UDP Client | Error encoding config to JSON", "error", err.Error())
		return
	}

	_, err = conn.Write(body)
	if err != nil {
		utils.Log.Errorw("UDP Client | Error broadcasting config", "error", err.Error())
		return
	}
	utils.Log.Infow("UDP Client | Broadcasted vote", "vote", electAction.Name)
}

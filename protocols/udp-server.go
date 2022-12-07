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

var electionStarted = false
var votesReceived = 0
var electionBoard = make(map[string]int)
var winner = ""
var maxVotes = 0

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
		var udpAction domain.UdpAction
		err = json.Unmarshal(buffer, &udpAction)
		if err != nil {
			utils.Log.Errorw("UDP Server | Error decoding config JSON", "error", err.Error(), "info", string(buffer))
			continue
		}

		switch udpAction.Action {

		case constants.UDP_CONFIG:
			if udpAction.Name != config.NAME {
				utils.Log.Infow("UDP Server | Received config", "config", udpAction)

				go UpdateTcpConnection(udpAction)
			}
		case constants.UDP_ELECT:
			utils.Log.Infow("UDP Server | Received vote", "vote", udpAction.Name)
			if !electionStarted {
				electionStarted = true
				utils.Log.Infow("UDP Server | Started election")
				go waitForElectionEnd()
			}

			electionBoard[udpAction.Name]++
			votesReceived++

			if electionBoard[udpAction.Name] > maxVotes {
				winner = udpAction.Name
				maxVotes = electionBoard[udpAction.Name]
			}
		}
	}
}

func waitForElectionEnd() {
	time.Sleep(constants.UDP_ELECTION_TIME)
	if votesReceived < len(datastore.Connections) {
		utils.Log.Warnw("UDP Server ELECTION | Election not enough votes", "votes", votesReceived, "board", electionBoard)
		go startElection()
		return
	}

	utils.Log.Infow("UDP Server ELECTION | Election finished", "winner", winner, "board", electionBoard)
	if winner == config.NAME {
		datastore.IsLeader = true
		NewLeaderChan <- true
	} else {
		UpdateTcpConnection(domain.UdpAction{
			Name:     winner,
			Port:     datastore.Connections[winner].Port,
			IsLeader: true,
		})

	}
	datastore.HasLeader = true

	electionStarted = false
	votesReceived = 0
	electionBoard = make(map[string]int)
	winner = ""
	maxVotes = 0
}

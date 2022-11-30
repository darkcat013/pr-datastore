package protocols

import (
	"encoding/json"
	"net"
	"time"

	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

func TcpInsert(id, value string) {
	action := domain.Action{
		Action: constants.TCP_INSERT,
		Value:  value,
		Id:     id,
	}

	tcpBroadcast(action)
}

func TcpUpdate(id, value string) {
	action := domain.Action{
		Action: constants.TCP_UPDATE,
		Value:  value,
		Id:     id,
	}

	tcpBroadcast(action)
}

func TcpDelete(id string) {
	action := domain.Action{
		Action: constants.TCP_DELETE,
		Id:     id,
	}

	tcpBroadcast(action)
}

func tcpBroadcast(action domain.Action) {
	utils.Log.Infow("TCP Client | Broadcast action", "action", action)

	data, err := json.Marshal(action)
	if err != nil {
		utils.Log.Errorw("TCP Client | Failed to encode to JSON", "error", err.Error(), "data", action)
		return
	}

	for _, connection := range datastore.Connections {
		_, err := connection.Conn.Write(data)
		if err != nil {
			utils.Log.Errorw("TCP Client | Failed to send data to address", "error", err.Error(), "address", connection.Conn.LocalAddr().String())
			continue
		}
	}
}

func UpdateTcpConnection(config domain.ConfigUdp) {

	utils.Log.Info("TCP Client | Update connection config for ", config.Name)
	if _, ok := datastore.Connections[config.Name]; ok {
		temp := datastore.Connections[config.Name]
		temp.LastUpdated = time.Now().Unix()
		datastore.Connections[config.Name] = temp
		return
	}

	addr := config.Name + config.Port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.Log.Error("TCP Client | Could not connect to address ", addr)
		return
	}

	datastore.Connections[config.Name] = datastore.Connection{
		Port:        config.Port,
		LastUpdated: time.Now().Unix(),
		Conn:        conn,
	}
}

package protocols

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/darkcat013/pr-datastore/config"
	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

var distributionMutex sync.Mutex

func TcpInsert(value []byte) string {
	newId := utils.GetNewId(value)
	action := domain.TcpAction{
		Action: constants.TCP_INSERT,
		Value:  value,
		Id:     newId,
	}

	go tcpDistributeData(action)
	return newId
}

func TcpUpdate(id string, value []byte) {
	action := domain.TcpAction{
		Action: constants.TCP_UPDATE,
		Value:  value,
		Id:     id,
	}

	tcpBroadcast(action)
}

func TcpDelete(id string) {
	action := domain.TcpAction{
		Action: constants.TCP_DELETE,
		Id:     id,
	}

	tcpBroadcast(action)
}

func tcpBroadcast(action domain.TcpAction) {
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

func tcpDistributeData(action domain.TcpAction) {
	utils.Log.Infow("TCP Client | Distribute data action", "action", action)

	data, err := json.Marshal(action)
	if err != nil {
		utils.Log.Errorw("TCP Client | Failed to encode to JSON", "error", err.Error(), "data", action)
		return
	}

	distributionMutex.Lock()

	distributeTo := make([]string, 0)
	currentKeys := datastore.GetAllKeys()
	maxFileAmount := len(currentKeys)
	distributedTo := 0
	distributeGoal := ((len(datastore.Connections) + 1) / 2) + 1
	if maxFileAmount == 0 {
		datastore.InsertAtId(action.Id, action.Value)
		currentKeys = datastore.GetAllKeys()
		maxFileAmount = len(currentKeys)
		distributedTo++
	}

	for connName, connection := range datastore.Connections {
		if len(connection.DataIds) < maxFileAmount {
			distributeTo = append(distributeTo, connName)
			connection.DataIds[action.Id] = true
			distributedTo++
		} else if maxFileAmount < len(connection.DataIds) {
			maxFileAmount = len(connection.DataIds)
			datastore.InsertAtId(action.Id, action.Value)
			distributedTo++
		}
	}

	if distributedTo < distributeGoal {
		if !utils.Contains(distributeTo, config.NAME) {
			datastore.InsertAtId(action.Id, action.Value)
			distributedTo++
		}
		for connName, connection := range datastore.Connections {
			if !utils.Contains(distributeTo, connName) && distributedTo < distributeGoal {
				distributeTo = append(distributeTo, connName)
				connection.DataIds[action.Id] = true
				distributedTo++
			}
		}
	}
	distributionMutex.Unlock()

	for i := 0; i < len(distributeTo); i++ {
		_, err := datastore.Connections[distributeTo[i]].Conn.Write(data)
		if err != nil {
			utils.Log.Errorw("TCP Client | Failed to send data to address", "error", err.Error(), "address", datastore.Connections[distributeTo[i]].Conn.LocalAddr().String())
			continue
		}
	}
}

func UpdateTcpConnection(config domain.UdpAction) {

	utils.Log.Info("TCP Client | Update connection config for ", config.Name)

	var conn net.Conn
	var err error

	if val, ok := datastore.Connections[config.Name]; ok {
		conn = val.Conn
	} else {
		addr := config.Name + config.Port
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			utils.Log.Error("TCP Client | Could not connect to address ", addr)
			return
		}
	}

	datastore.Connections[config.Name] = datastore.Connection{
		Port:        config.Port,
		LastUpdated: time.Now().Unix(),
		Conn:        conn,
		IsLeader:    config.IsLeader,
		DataIds:     config.DataIds,
	}
}

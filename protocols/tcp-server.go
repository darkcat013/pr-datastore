package protocols

import (
	"encoding/json"
	"net"

	"github.com/darkcat013/pr-datastore/config"
	"github.com/darkcat013/pr-datastore/constants"
	"github.com/darkcat013/pr-datastore/datastore"
	"github.com/darkcat013/pr-datastore/domain"
	"github.com/darkcat013/pr-datastore/utils"
)

func StartTcpServer() {

	l, err := net.Listen("tcp", config.TCP_PORT)
	if err != nil {
		utils.Log.Errorw("TCP Server | Could not start TCP", "error", err.Error())
		return
	}

	utils.Log.Info("TCP Server | Start listening on port ", config.TCP_PORT)

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			utils.Log.Errorw("TCP Server | Error accepting connection", "error", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			utils.Log.Errorw("TCP Server | Error reading from connection", "error", err.Error())
			return
		}

		buffer = buffer[:n]
		var action domain.TcpAction
		err = json.Unmarshal(buffer, &action)
		if err != nil {
			utils.Log.Errorw("TCP Server | Error decoding config JSON", "error", err.Error(), "info", string(buffer))
			continue
		}

		utils.Log.Infow("TCP Server | Received action", "action", action)

		go handleAction(action)
	}
}

func handleAction(action domain.TcpAction) {
	switch action.Action {
	case constants.TCP_INSERT:
		err := datastore.InsertAtId(action.Id, action.Value)
		if err != nil {
			utils.Log.Errorw("TCP Server | Could not execute TCP_INSERT", "action", action, "error", err)
			return
		}
	case constants.TCP_UPDATE:
		err := datastore.Update(action.Id, action.Value)
		if err != nil {
			utils.Log.Errorw("TCP Server | Could not execute TCP_UPDATE", "action", action, "error", err)
			return
		}
	case constants.TCP_DELETE:
		err := datastore.Delete(action.Id)
		if err != nil {
			utils.Log.Errorw("TCP Server | Could not execute TCP_DELETE", "action", action, "error", err)
			return
		}
	}
}

package datastore

import "net"

type Connection struct {
	Port        string
	Conn        net.Conn
	LastUpdated int64
}

var Connections = make(map[string]Connection)

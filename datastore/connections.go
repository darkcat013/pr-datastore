package datastore

import "net"

type Connection struct {
	Port        string
	Conn        net.Conn
	LastUpdated int64
	IsLeader    bool
	DataIds     map[string]bool
}

var Connections = make(map[string]Connection)

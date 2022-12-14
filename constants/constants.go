package constants

import "time"

const HTTP_PORT = ":8080"

const UDP_BROADCAST_IP = "255.255.255.255"
const UDP_BROADCAST_PORT = ":8075"
const UDP_BROADCAST_DELAY = 5 * time.Second
const UDP_ELECTION_TIME = 2 * time.Second

const CONNECTION_ALIVE_TIME = 8

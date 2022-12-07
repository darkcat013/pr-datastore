package datastore

import "github.com/darkcat013/pr-datastore/config"

var IsLeader bool = config.NAME == "datastore-1"
var HasLeader bool = true

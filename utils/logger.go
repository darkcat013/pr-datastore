package utils

import (
	"github.com/darkcat013/pr-datastore/config"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitializeLogger() {
	var tempLog *zap.Logger
	if config.LOGS_ENABLED {
		tempLog, _ = zap.NewDevelopment()
	} else {
		tempLog = zap.NewNop()
	}

	Log = tempLog.Sugar()
	Log.Sync()
}

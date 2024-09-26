package config

import (
	"github.com/HomesNZ/go-common/env"
	"github.com/HomesNZ/go-common/logger"
)

func InitLogger() {
	logger.Init(
		logger.Level(env.GetString("LOG_LEVEL", "info")),
	)
}

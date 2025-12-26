package main

import (
	"github.com/tuanta7/k6noz/services/internal/ingestion"
	"github.com/tuanta7/k6noz/services/pkg/serverx"
	"github.com/tuanta7/k6noz/services/pkg/zapx"
	"go.uber.org/zap"
)

func main() {
	cfg, err := ingestion.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := zapx.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	handler := ingestion.NewHandler(logger)
	server := ingestion.NewServer(cfg, handler)

	logger.Info("starting server", zap.String("address", cfg.BindAddress))
	if err = serverx.RunServer(server); err != nil {
		logger.Error("failed to run server", zap.Error(err))
		return
	}
}

package zapx

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	*zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"

	zl, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	return &ZapLogger{zl}, nil
}

func (zl *ZapLogger) Close() error {
	return zl.Logger.Sync()
}

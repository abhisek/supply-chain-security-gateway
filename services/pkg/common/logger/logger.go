package logger

import "go.uber.org/zap"

var (
	defaultLogger = zap.NewNop()
	sugarLogger   = defaultLogger.Sugar()
)

// Init logger for a service
// There are two type of logging config:
//		1. Service specific
//		2. Common config
// Common config helps in logging consistency across all services
// in an environment
func Init(svc string) {
	l, err := zapBuild(zapConfig())
	if err != nil {
		panic("Failed to build logger")
	}

	defaultLogger = l.With(zap.String("service", svc))
	sugarLogger = defaultLogger.Sugar()
}

func zapConfig() zap.Config {
	return zap.NewProductionConfig()
}

func zapBuild(config zap.Config) (*zap.Logger, error) {
	return config.Build()
}

func Infof(msg string, args ...any) {
	sugarLogger.Infof(msg, args...)
}

func Warnf(msg string, args ...any) {
	sugarLogger.Infof(msg, args...)
}

func Errorf(msg string, args ...any) {
	sugarLogger.Errorf(msg, args...)
}

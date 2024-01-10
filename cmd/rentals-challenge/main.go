package main

import (
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	serviceID = "rental-challenge"
)

var cli struct {
	Debug      int    `kong:"short='d',env='DEBUG',default=0,help='Run in debug mode'"`
	HTTPPort   int    `kong:"short='t',env='HTTP_PORT',default='59191',help='HTTP server port'"`
	DBHost     string `kong:"short='h',env='DB_HOST',default='127.0.0.1',help='DB server host'"`
	DBPort     int    `kong:"short='r',env='DB_PORT',default='5434',help='DB server port'"`
	DBName     string `kong:"short='n',env='DB_NAME',default='testingwithrentals',help='DB name'"`
	DBUsername string `kong:"short='u',env='DB_USERNAME',default='root',help='DB username'"`
	DBPassword string `kong:"short='p',env='DB_PASSWORD',default='root',help='DB password'"`
}

func main() {
	kong.Parse(&cli, kong.Name(serviceID), kong.Description("Rental Service Code Challenge"), kong.UsageOnError())
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.TimeKey = "time"
	logCfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	if cli.Debug != 0 {
		logCfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := logCfg.Build()
	if err != nil {
		logger.Fatal("Failed to initialize logger. Exiting " + err.Error())
	}
	logger.Info("Starting rental service")
}

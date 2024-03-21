package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	initLogger()
	initViper()
	server := InitWebServer()
	server.Run(":8080")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func initViper() error {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return nil
}

func initViperRemote() error {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	return nil
}

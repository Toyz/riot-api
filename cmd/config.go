package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	viper.SetDefault("log_level", "info")

	viper.SetConfigFile("config.json")
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("no config file is found... Falling back to using Environment Variables")
	}

	level, err := log.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Warn("invalid log level... Falling back to info")
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

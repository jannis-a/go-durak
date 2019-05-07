package app

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/jannis-a/go-durak/utils"
)

type Config struct {
	BIND string
	KEY  string
	DB   string
}

func InitConfig() {
	viper.AddConfigPath(utils.GetPackagePath())
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
}

func NewConfig() *Config {
	InitConfig()

	config := new(Config)
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Panic(err)
	}

	return config
}

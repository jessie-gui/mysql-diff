package helper

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/jessie-gui/mysql-diff/model"
	"github.com/spf13/viper"
)

// NewConfig /**
func NewConfig() *model.Config {
	v := viper.New()
	c := &model.Config{}

	v.SetConfigFile("config/config.yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("read config/config.yaml file failed: ", err)
	}

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file is updated: %v", e.Name)

		if err := viper.Unmarshal(c); err != nil {
			log.Printf("config file OnConfigChange Unmarshal failed: %v", err)
		}
	})

	if err := v.Unmarshal(c); err != nil {
		log.Fatal("config file Unmarshal failed: ", err)
	}

	return c
}

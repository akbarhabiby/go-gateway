package config

import (
	"github.com/labstack/gommon/color"
	"github.com/spf13/viper"
)

var (
	Config *viper.Viper
)

func init() {
	Config = viper.New()
	Config.SetConfigFile(".env")

	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
	color.Println(color.Green("⇨ config loaded"))
}

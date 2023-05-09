package config

import (
	"github.com/labstack/gommon/color"
	"github.com/spf13/viper"
)

var (
	// v viper global variable
	Config *viper.Viper
	// env used environment
)

func init() {
	// initializing viper
	Config = viper.New()
	Config.SetConfigFile(".env")

	if err := Config.ReadInConfig(); err != nil {
		panic(err)
	}
	color.Println(color.Green("⇨ config loaded"))
}

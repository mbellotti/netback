package main

import (
	"fmt"

	"github.com/mbellotti/netback/cmd"
	_ "github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Configuration not found, please run `netback config` to set up")
	} else {
		viper.WatchConfig()
	}
	cmd.Execute()
}

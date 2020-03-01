package main

import (
	"fmt"
	"os"

	"github.com/mbellotti/netback/cmd"
	"github.com/mitchellh/go-homedir"
	_ "github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	h, err := homedir.Dir()
	configpath := fmt.Sprintf("%s/netback", h)
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		err := os.Mkdir(configpath, 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
	filepath := fmt.Sprintf("%s/config.yaml", configpath)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		_, err = os.Create(filepath)
		if err != nil {
			fmt.Println(err)
		}
	}

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configpath) // path to look for the config file in

	fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Configuration not found, please run `netback config` to set up")
	} else {
		viper.WatchConfig()
	}
	cmd.Execute()
}

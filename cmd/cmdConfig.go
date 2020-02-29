package cmd

import (
	"errors"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string

	validate = func(input string) error {
		if input == "" {
			return errors.New("Invalid input")
		}
		return nil
	}
	cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "Configure Netback",
		Long:  `Configure Twitter API settings`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if ConsumerKey != "" {
				viper.Set("ConsumerKey", ConsumerKey)
			} else if ConsumerSecret != "" {
				viper.Set("ConsumerSecret", ConsumerSecret)
			} else if Token != "" {
				viper.Set("Token", Token)
			} else if TokenSecret != "" {
				viper.Set("TokenSecret", TokenSecret)
			} else {
				prompt := promptui.Prompt{
					Label:    "Twiter Consumer Key:",
					Validate: validate,
				}

				ConsumerKey, err := prompt.Run()
				if err != nil {
					log.Fatalf(err.Error())
				} else {
					viper.Set("ConsumerKey", ConsumerKey)
				}

				prompt = promptui.Prompt{
					Label:    "Twiter Consumer Secret:",
					Validate: validate,
				}

				ConsumerSecret, err := prompt.Run()
				if err != nil {
					log.Fatalf(err.Error())
				} else {
					viper.Set("ConsumerSecret", ConsumerSecret)
				}

				prompt = promptui.Prompt{
					Label:    "Twiter Token:",
					Validate: validate,
				}

				Token, err := prompt.Run()
				if err != nil {
					log.Fatalf(err.Error())
				} else {
					viper.Set("Token", Token)
				}

				prompt = promptui.Prompt{
					Label:    "Twiter Token Secret:",
					Validate: validate,
				}

				TokenSecret, err := prompt.Run()
				if err != nil {
					log.Fatalf(err.Error())
				} else {
					viper.Set("TokenSecret", TokenSecret)
				}
			}
			viper.WriteConfig()
		},
	}
)

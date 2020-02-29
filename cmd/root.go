package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netback",
	Short: "Are suspicious accounts connected to yours?",
	Long: `
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
 _   _      _   _                _           ______
| \ | |    | | | |              | |       .-'      \.
|  \| | ___| |_| |__   __ _  ___| | __   /           \
| .   |/ _ \ __|  _ \ / _  |/ __| |/ /  |-.           |
| |\  |  __/ |_| |_) | (_| | (__|   <   |  \          |
|_| \_|\___|\__|_.__/ \__,_|\___|_|\_\  [__|__________|_______

$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$
											  
Netback is a command line tool to examine your Twitter connections and identify accounts with suspicious activity.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	cmdConfig.Flags().StringVarP(&ConsumerKey, "consumer-key", "k", "", "Twitter consumer key")
	cmdConfig.Flags().StringVarP(&ConsumerSecret, "consumer-secret", "c", "", "Twitter consumer secret")
	cmdConfig.Flags().StringVarP(&Token, "token", "t", "", "Twitter auth token")
	cmdConfig.Flags().StringVarP(&TokenSecret, "token-secret", "s", "", "Twitter token secret")
	cmdFriends.Flags().BoolP("verbose", "v", false, "Print results for all accounts examined")
	cmdFollowers.Flags().BoolP("verbose", "v", false, "Print results for all accounts examined")
	cmdSelf.Flags().BoolP("verbose", "v", false, "Print results for profile analysis")
	cmdProfile.Flags().BoolP("low-activity", "l", false, "Configure a low activity bot profile")
	cmdProfile.Flags().BoolP("high-activity", "j", false, "Configure a high activity bot profile")

	rootCmd.AddCommand(cmdFollowers)
	rootCmd.AddCommand(cmdFriends)
	rootCmd.AddCommand(cmdSelf)
	rootCmd.AddCommand(cmdConfig)
	rootCmd.AddCommand(cmdProfile)
}

func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	// Search config in home directory with name ".cobra" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".cobra")
	// }

	// viper.AutomaticEnv()

	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}

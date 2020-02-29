package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdSelf = &cobra.Command{
	Use:   "self",
	Short: "Would your account look suspicious?",
	Long:  `Check how your own account compares`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := strconv.ParseBool(cmd.Flag("verbose").Value.String())
		profileTest := viper.GetString("FollowerDirection")
		if profileTest == "" {
			log.Fatalf("Please define characteristics of suspicious accounts by running the command `netback profile`")
		} else {
			tw := NewTwitterConnection()
			for !tw.RateLimit() {
				log.Println("Twitter rate limt reached. Sleeping...")
				time.Sleep(15 * time.Minute)
			}
			if viper.GetInt("PPDWeight") == 0 && viper.GetInt("FollowerWeight") == 0 && viper.GetInt("PerOrgWeight") == 0 && viper.GetInt("ProfilePicWeight") == 0 && viper.GetInt("LowPostsWeight") == 0 {
				log.Fatalf("Profile contains no criteria with which to assess accounts.")
			}

			tweets, _ := tw.FetchTimelineData(0)
			score, _ := BotScore(tweets, v)
			if v {
				fmt.Printf("Account: %s Score: %d\n", tweets[0].User.ScreenName, score)
				fmt.Println("--------------------------")
			}
			if score >= viper.GetInt("Score") {
				fmt.Println("Your account would be considered suspicious with these settings. Run `netback profile` to reset the profile with different settings.")
			} else {
				fmt.Println("Your account would not be considered suspicious with these settings.")
			}
		}

	},
}

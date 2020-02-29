package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdFollowers = &cobra.Command{
	Use:   "followers",
	Short: "Search amoung your followers",
	Long:  `Examine accounts that follow yours for suspicious behavior. Takes no arguments`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := strconv.ParseBool(cmd.Flag("verbose").Value.String())
		profileTest := viper.GetString("FollowerDirection")
		if profileTest == "" {
			log.Fatalln("Please define characteristics of suspicious accounts by running the command `netback profile`")
		} else {
			tw := NewTwitterConnection()
			for !tw.RateLimit() {
				log.Println("Twitter rate limt reached. Sleeping...")
				time.Sleep(15 * time.Minute)
			}
			if viper.GetInt("PPDWeight") == 0 && viper.GetInt("FollowerWeight") == 0 && viper.GetInt("PerOrgWeight") == 0 && viper.GetInt("ProfilePicWeight") == 0 && viper.GetInt("LowPostsWeight") == 0 {
				log.Fatalln("Profile contains no criteria with which to assess accounts.")
				return
			}
			u, _ := tw.FetchTimelineData(0)
			NumOfFollowers := u[0].User.FollowersCount
			fmt.Printf("Beginning analysis: %d accounts to look at\n", NumOfFollowers)

			var bar *pb.ProgressBar
			if !v {
				bar = pb.StartNew(NumOfFollowers)
			}

			cursor := viper.GetInt64("FollowersCursor")
			if cursor == 0 { //Zero value is end of list
				cursor = -1
			}
			for cursor != int64(0) {
				users, res, c := tw.FetchFollowerList(cursor)
				cursor = c

				if RateLimitGood(res) {
					fmt.Println("Twitter rate limt reached. Sleeping for 15 minutes.")
					viper.Set("FollowersCursor", cursor)
					viper.WriteConfig()
					time.Sleep(15 * time.Minute)
				}
				for _, f := range users.Users {
					if !v {
						bar.Increment()
					}
					tweets, res := tw.FetchTimelineData(f.ID)

					if RateLimitGood(res) {
						fmt.Println("Twitter rate limt reached. Sleeping for 15 minutes.")
						viper.Set("FollowersCursor", cursor)
						viper.WriteConfig()
						time.Sleep(15 * time.Minute)
					}
					score, username := BotScore(tweets, v)
					if score >= viper.GetInt("Score") || v {
						fmt.Printf("Account: %s Score: %d\n", username, score)
						fmt.Println("--------------------------")
					}
				}
			}
			viper.Set("FollowersCursor", -1)
			viper.WriteConfig()
			if !v {
				bar.Finish()
			}
		}
	},
}

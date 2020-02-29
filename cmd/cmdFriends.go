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

var cmdFriends = &cobra.Command{
	Use:   "friends",
	Short: "Search amoung your friends",
	Long:  `Examine accounts you are following for suspicious behavior`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		v, _ := strconv.ParseBool(cmd.Flag("verbose").Value.String())
		profileTest := viper.GetString("FollowerDirection")
		if profileTest == "" {
			log.Fatalf("Please define characteristics of suspicious accounts by running the command `netback profile`")
		} else {
			tw := NewTwitterConnection()
			for !tw.RateLimit() {
				fmt.Println("Twitter rate limt reached. Sleeping...")
				time.Sleep(15 * time.Minute)
			}
			if viper.GetInt("PPDWeight") == 0 && viper.GetInt("FollowerWeight") == 0 && viper.GetInt("PerOrgWeight") == 0 && viper.GetInt("ProfilePicWeight") == 0 && viper.GetInt("LowPostsWeight") == 0 {
				log.Fatalf("Profile contains no criteria with which to assess accounts.")
				return
			}
			u, _ := tw.FetchTimelineData(0)
			NumOfFriends := u[0].User.FriendsCount
			fmt.Printf("Beginning analysis: %d accounts to look at\n", NumOfFriends)

			var bar *pb.ProgressBar
			if !v {
				bar = pb.StartNew(NumOfFriends)
			}

			cursor := viper.GetInt64("FriendsCursor")
			if cursor == 0 { //Zero value is end of list
				cursor = -1
			}
			for cursor != int64(0) {
				users, res, c := tw.FetchFriendList(cursor)
				cursor = c

				if RateLimitGood(res) {
					fmt.Println("Twitter rate limt reached. Sleeping for 15 minutes.")
					viper.Set("FriendsCursor", cursor)
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
						viper.Set("FriendsCursor", cursor)
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
			viper.Set("FriendsCursor", -1)
			viper.WriteConfig()
			if !v {
				bar.Finish()
			}
		}
	},
}

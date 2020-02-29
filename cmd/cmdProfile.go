package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	PostsPerDay          int
	PPDWeight            int
	FollowerThresholdRaw string
	FollowerThreshold    int64
	FollowerDirection    string
	FollowerWeight       int
	PercentOriginal      int
	PerOrgWeight         int
	ProfilePicWeight     int
	LowPosts             int
	LowPostsWeight       int
	Score                int

	cmdProfile = &cobra.Command{
		Use:   "profile",
		Short: "Configure Netback's bot profile",
		Long:  `Configure what Netback thinks a bot is`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			l, _ := strconv.ParseBool(cmd.Flag("low-activity").Value.String())
			j, _ := strconv.ParseBool(cmd.Flag("high-activity").Value.String())
			a := false
			var prompt promptui.Prompt
			if !l && !j {
				a = true
			}
			if a || j {
				if j {
					viper.Set("ProfilePicWeight", 0)
					viper.Set("LowPostsWeight", 0)
				}
				prompt = promptui.Prompt{
					Label: "Posts per day:",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" || string(input[0]) == "." {
							return errors.New("Invalid number")
						}

						if v < 0 {
							return errors.New("Invalid number")
						}
						return nil
					},
				}

				PostsPerDay, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("PostsPerDay", PostsPerDay)
				}

				prompt = promptui.Prompt{
					Label: "Weight to give Posts Per Day (0-3):",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if v < 0 || v > 3 {
							return errors.New("Invalid weight. Must be between 0-3")
						}
						return nil
					},
				}

				PPDWeight, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("PPDWeight", PPDWeight)
				}

				prompt = promptui.Prompt{
					Label: "Percentage of post that are not original (0-100):",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if v < 0 || v > 100 {
							return errors.New("Invalid input. Must be between 0-100")
						}
						return nil
					},
				}

				PercentOriginal, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("PercentOriginal", PercentOriginal)
				}

				prompt = promptui.Prompt{
					Label: "Weight to give percentage of original posts (0-3):",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if v < 0 || v > 3 {
							return errors.New("Invalid weight. Must be between 0-3")
						}
						return nil
					},
				}

				PerOrgWeight, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("PerOrgWeight", PerOrgWeight)
				}
			}

			prompt = promptui.Prompt{
				Label: "Max or min number of followers (+0, -1000):",
				Validate: func(input string) error {
					_, err := strconv.ParseFloat(input, 64)
					if err != nil {
						return errors.New("Invalid number")
					}

					if string(input[0]) != "-" && string(input[0]) != "+" {
						return errors.New(string(input[0]) + " Invalid input, must start with a + or a -")
					}

					return nil
				},
			}

			FollowerThresholdRaw, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
			} else {
				if FollowerThresholdRaw != "" {
					FollowerDirection = string(FollowerThresholdRaw[0])
					FollowerThreshold, _ = strconv.ParseInt(FollowerThresholdRaw[1:], 10, 64)
					viper.Set("FollowerDirection", FollowerDirection)
					viper.Set("FollowerThreshold", FollowerThreshold)
				}
			}

			prompt = promptui.Prompt{
				Label: "Weight to give Follower counts (0-3):",
				Validate: func(input string) error {
					v, err := strconv.ParseFloat(input, 64)
					if err != nil {
						return errors.New("Invalid number")
					}

					if string(input[0]) == "-" {
						return errors.New("Invalid number")
					}

					if v < 0 || v > 3 {
						return errors.New("Invalid weight. Must be between 0-3")
					}
					return nil
				},
			}

			FollowerWeight, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
			} else {
				viper.Set("FollowerWeight", FollowerWeight)
			}

			if a || l {
				if l {
					viper.Set("PPDWeight", 0)
					viper.Set("PerOrgWeight", 0)
				}
				prompt = promptui.Prompt{
					Label: "Weight to give the lack of profile picture (0-3):",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if v < 0 || v > 3 {
							return errors.New("Invalid weight. Must be between 0-3")
						}
						return nil
					},
				}

				ProfilePicWeight, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("ProfilePicWeight", ProfilePicWeight)
				}

				prompt = promptui.Prompt{
					Label: "A suspicious account has fewer than how many posts?",
					Validate: func(input string) error {
						_, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if strings.Contains(string(input), ".") {
							return errors.New("Invalid number")
						}

						return nil
					},
				}

				LowPosts, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("LowPosts", LowPosts)
				}

				prompt = promptui.Prompt{
					Label: "Weight to give a low post count (0-3):",
					Validate: func(input string) error {
						v, err := strconv.ParseFloat(input, 64)
						if err != nil {
							return errors.New("Invalid number")
						}

						if string(input[0]) == "-" {
							return errors.New("Invalid number")
						}

						if v < 0 || v > 3 {
							return errors.New("Invalid weight. Must be between 0-3")
						}
						return nil
					},
				}

				LowPostsWeight, err := prompt.Run()
				if err != nil {
					fmt.Println(err)
				} else {
					viper.Set("LowPostsWeight", LowPostsWeight)
				}
			}

			prompt = promptui.Prompt{
				Label: "Total score for bots (1-15):",
				Validate: func(input string) error {
					v, err := strconv.ParseFloat(input, 64)
					if err != nil {
						return errors.New("Invalid number")
					}

					if string(input[0]) == "-" {
						return errors.New("Invalid number")
					}

					if v < 1 || v > 16 {
						return errors.New("Invalid score. Must be between 1-15")
					}
					return nil
				},
			}

			Score, err := prompt.Run()
			if err != nil {
				fmt.Println(err)
			} else {
				viper.Set("Score", Score)
			}

			viper.WriteConfig()
		},
	}
)

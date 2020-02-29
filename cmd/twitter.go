package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

type Twitter struct {
	client *twitter.Client
}

func NewTwitterConnection() *Twitter {
	config := oauth1.NewConfig(viper.GetString("ConsumerKey"), viper.GetString("ConsumerSecret"))
	token := oauth1.NewToken(viper.GetString("Token"), viper.GetString("TokenSecret"))
	httpClient := config.Client(oauth1.NoContext, token)
	t := Twitter{client: twitter.NewClient(httpClient)}
	return &t
}

func (t Twitter) RateLimit() bool {
	rl, _, err := t.client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{"statuses"}})
	if err != nil {
		return false
	}
	if rl.Resources.Statuses["/statuses/user_timeline"].Remaining > 0 {
		return true
	}
	return false
}

func (t Twitter) FetchFriendList(cursor int64) (*twitter.Friends, *http.Response, int64) {
	friends, res, err := t.client.Friends.List(&twitter.FriendListParams{UserID: 0, Cursor: cursor})
	if err != nil {
		fmt.Println(err.Error())
	}
	return friends, res, friends.NextCursor
}

func (t Twitter) FetchFollowerList(cursor int64) (*twitter.Followers, *http.Response, int64) {
	followers, res, err := t.client.Followers.List(&twitter.FollowerListParams{UserID: 0, Cursor: cursor})
	if err != nil {
		fmt.Println(err)
	}
	return followers, res, followers.NextCursor
}

func (t Twitter) FetchTimelineData(UserID int64) ([]twitter.Tweet, *http.Response) {
	tweets, res, err := t.client.Timelines.UserTimeline(&twitter.UserTimelineParams{UserID: UserID, IncludeRetweets: twitter.Bool(true)})
	if err != nil {
		fmt.Println(err.Error())
	}
	return tweets, res
}

func (t Twitter) PPD(tweets []twitter.Tweet) float64 {
	StatusesCount := tweets[0].User.StatusesCount
	tm, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", tweets[0].User.CreatedAt)
	days := daysBetween(time.Now(), tm)
	if days == 0 {
		fmt.Println("This account is brand new, draw your own conclusions.")
	}
	return (float64(StatusesCount) / float64(days))
}

func (t Twitter) Followers(tweets []twitter.Tweet) int {
	return tweets[0].User.FollowersCount
}

func (t Twitter) PercentOriginal(tweets []twitter.Tweet) int {
	count := 0
	retweet := 0
	for _, f := range tweets {
		count = count + 1
		if f.RetweetedStatus != nil {
			retweet = retweet + 1
		}
	}
	if retweet == 0 {
		return 0
	}
	return int((float64(retweet) / float64(count)) * 100)
}

func (t Twitter) DefaultProfilePic(tweets []twitter.Tweet) bool {
	return tweets[0].User.DefaultProfileImage
}

func (t Twitter) LowPosts(tweets []twitter.Tweet) int {
	return tweets[0].User.StatusesCount
}

func RateLimitGood(res *http.Response) bool {
	if res.Header["X-Rate-Limit-Remaining"] != nil {
		return res.Header["X-Rate-Limit-Remaining"][0] == "0"
	}

	if res.Header["X-App-Rate-Limit-Remaining"] != nil {
		return res.Header["X-App-Rate-Limit-Remaining"][0] == "0"
	}
	return false
}

func BotScore(tweets []twitter.Tweet, v bool) (int, string) {
	if len(tweets) == 0 {
		return 0, "" //Typically means account is deleted or suspended
	}
	tw := Twitter{}
	score := 0

	if viper.GetInt("PPDWeight") > 0 {
		PostsPerDay := tw.PPD(tweets)
		if v {
			fmt.Printf("%s: posts per day observed: %f allowed: %d \n", tweets[0].User.ScreenName, PostsPerDay, viper.GetInt("PostsPerDay"))
		}
		if PostsPerDay >= float64(viper.GetInt("PostsPerDay")) {
			score = score + viper.GetInt("PPDWeight")
		}
	}

	if viper.GetInt("FollowerWeight") > 0 {
		FollowerCount := tw.Followers(tweets)
		if viper.GetString("FollowerDirection") == "+" {
			if v {
				fmt.Printf("%s: number of followers observed: %d allowed: under %d \n", tweets[0].User.ScreenName, FollowerCount, viper.GetInt("FollowerThreshold"))
			}
			if FollowerCount >= viper.GetInt("FollowerThreshold") {
				score = score + viper.GetInt("FollowerWeight")
			}
		}
		if viper.GetString("FollowerDirection") == "-" {
			if v {
				fmt.Printf("%s: number of followers observed: %d allowed: over %d \n", tweets[0].User.ScreenName, FollowerCount, viper.GetInt("FollowerThreshold"))
			}
			if FollowerCount <= viper.GetInt("FollowerThreshold") {
				score = score + viper.GetInt("FollowerWeight")
			}
		}
	}
	if viper.GetInt("PerOrgWeight") > 0 {
		originals := tw.PercentOriginal(tweets)
		if v {
			fmt.Printf("%s: percentage of tweets with no original content observed: %d allowed: under %d \n", tweets[0].User.ScreenName, originals, viper.GetInt("PercentOriginal"))
		}

		if originals >= viper.GetInt("PercentOriginal") {
			score = score + viper.GetInt("PerOrgWeight")
		}
	}

	if viper.GetInt("ProfilePicWeight") > 0 {
		pic := tw.DefaultProfilePic(tweets)
		if v {
			if pic {
				fmt.Printf("%s: using the default profile pic\n", tweets[0].User.ScreenName)
			} else {
				fmt.Printf("%s: using a custom profile pic\n", tweets[0].User.ScreenName)
			}
		}
		if pic {
			score = score + viper.GetInt("ProfilePicWeight")
		}
	}

	if viper.GetInt("LowPostsWeight") > 0 {
		posts := tw.LowPosts(tweets)
		if v {
			fmt.Printf("%s: number of total tweets observed: %d allowed: over %d \n", tweets[0].User.ScreenName, posts, viper.GetInt("LowPosts"))
		}
		if posts <= viper.GetInt("LowPosts") {
			score = score + viper.GetInt("LowPostsWeight")
		}
	}
	return score, tweets[0].User.ScreenName
}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

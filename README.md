# Netback
A tool for finding suspicious Twitter accounts that are connected to yours

## Installation
Netback is written in Go, so follow [these instructions](https://golang.org/doc/install) to make sure Go is installed on your machine. Then just run the following command:

```
go get github.com/mbellotti/netback
go install github.com/mbellotti/netback
```

You can also download it from source and execute the program directly from its directory if you prefer. Commands in that case would be `go run main.go [command] -flags` instead of `netback [command] -flags`

## Configuration
Netback needs Twitter API access through your account. The program does not need anything more than read permission on public Twitter data. You'll need a [developer account](https://developer.twitter.com/) on Twitter. Create a new app and copy the API key, API Secret, access token and token secret. Then run `netback config` and enter in that information.

## Setting up a Profile
The bot profile is the criteria Netback uses to determine whether an account looks suspicious or not. You can configure the profile to check for low activity accounts that might be fake followers or high activity accounts or some custom combination of features from both. Netback will ask you to set a weight for each criteria, assigning a weight of 0 effectively turns that criteria off.

The current possible criteria to use include:
- Posts per Day: calculated as total posts over days the account has existed.
- Number of Followers: with direction, meaning you can tell netback that anything over 100,000 followers (+100000) is suspicious or anything under 5 followers (-5) is suspicious.
- Percentage of Unoriginal Posts: Retweets without any original commentary.
- Profile Pic: Do they have one?
- Low Posts: Are their total number of posts too low?

Netback includes two flags `--low-activity` and `--high-activity` that act like short cuts, zeroing out criteria not relevant and only prompting for a subsection.

## Setting a Score
At the end of the profile wizard, Netback will prompt you to set a score. This determines the number of points an account needs to be labeled suspicious. For example if your profile gives percentage of unoriginal posts a weight of 3 and Number of Followers a weight of 2 then a logical bot score would be 5. With more criteria, setting up a bot score lower than the maximum value of all weights allows for some flexibility in how accounts are assessed.

## Testing your Profile
After building a profile and setting a score it's a good idea to test Netback's analysis against an account that you know to be legitimate-- your own. Run `netback self` will compare your own account against the profile you've built and provide feedback.

## Rate Limits
Twitter API only allows for a certain number of requests in a 15 minute period. When the rate limit is hit, Netback will alert the user, save it's current cursor position and sleep for fifteen minutes. Because the cursor position is saved, you can terminate the program while it's sleeping and restart it later on.

## Followers -vs- Friends
Twitter defines followers as people who follow you, and friends as people **you** follow. Netback will search for suspicious account in either list.

## Verbose Mode
`netback self` `netback friends` and `netback followers` can all be run in verbose mode which will output the observed values for every criteria with the allowed values set in the profile. Occassionally Netback will return a blank user. Typically this means the account has been suspended. It is still listed as a friend/follower but is otherwise inaccessible.

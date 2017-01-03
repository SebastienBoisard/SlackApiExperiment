package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"log"
)

func printUser(user slack.User) {
	fmt.Printf("   id=%s\n", user.ID)
	fmt.Printf("   name=%s\n", user.Name)
	fmt.Printf("   real name=%s\n", user.RealName)
	if user.IsBot == true {
		fmt.Printf("   is a bot\n")
	}
	if user.Deleted == true {
		fmt.Printf("   is a deactivated user\n")
	}
	// 'active'|'away'|''
	fmt.Printf("   presence=%s\n", user.Presence)
}

// type User struct {
//     ID                string      `json:"id"`
//     Name              string      `json:"name"`
//     Deleted           bool        `json:"deleted"`
//     Color             string      `json:"color"`
//     RealName          string      `json:"real_name"`
//     TZ                string      `json:"tz,omitempty"`
//     TZLabel           string      `json:"tz_label"`
//     TZOffset          int         `json:"tz_offset"`
//     Profile           UserProfile `json:"profile"`
//     IsBot             bool        `json:"is_bot"`
//     IsAdmin           bool        `json:"is_admin"`
//     IsOwner           bool        `json:"is_owner"`
//     IsPrimaryOwner    bool        `json:"is_primary_owner"`
//     IsRestricted      bool        `json:"is_restricted"`
//     IsUltraRestricted bool        `json:"is_ultra_restricted"`
//     Has2FA            bool        `json:"has_2fa"`
//     HasFiles          bool        `json:"has_files"`
//     Presence          string      `json:"presence"`
// }

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Error: no configuration file found")
		return
	}

	slackToken := viper.GetString("connection.token")

	slackHandler := slack.New(slackToken)

	// GetUsers returns the list of users (with their detailed information)
	// func (api *Client) GetUsers() ([]User, error)
	// Equivalent to (users.list method)[https://api.slack.com/methods/users.list]
	users, err := slackHandler.GetUsers()
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	for i, user := range users {
		fmt.Printf("user[%d]\n", i)
		printUser(user)
	}

   // Test PostMessage
	commands := map[string]string{
		"top":          "See the top rank of user rating by a provided number of top spots.",
		"bottom":       "See the bottom rank of user rating by a provided number of bottom spots.",
		"help":         "See the available bot commands.",
		"mean":         "See how the rating of the selected user looks like, comparing to the mean of all users.",
		"mean of":      "See how the rating of the selected user looks like, comparing to the mean of all users.",
		"top messages": "See the top ranking messages in the current channel."}

	fields := make([]slack.AttachmentField, 0)
	for k, v := range commands {
		fields = append(fields, slack.AttachmentField{
			Title: "<bot> " + k,
			Value: v,
		})
	}

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "Command List",
		Color:   "#A733FF",
		Text:    "some text",
		Fields:  fields,
	}

	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := slackHandler.PostMessage("#general", "Main text", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
}

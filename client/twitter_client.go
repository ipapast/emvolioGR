package client

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func AddDataToTweet(dataValue float64, textValue string) string {
	// bar_total is 15 and a new bar is every 6.67%
	dataToAdd := ""
	solid_bars_to_print := math.Round(dataValue / 6.67)
	empty_bars_to_print := 15 - int(solid_bars_to_print)

	dataToAdd += textValue
	for solid_bars_to_print > 0 {
		dataToAdd += "▓"
		solid_bars_to_print -= 1
	}
	for empty_bars_to_print > 0 {
		dataToAdd += "░"
		empty_bars_to_print -= 1
	}
	dataToAdd += " " + fmt.Sprint(dataValue) + "%\n\n"
	return dataToAdd
}

func SourceAndSendTweet(stringToTweet, language string) {
	now := time.Now()
	date := fmt.Sprint(now.Format("02/01/2006"))

	if language == "en" {
		stringToTweet += "As of " + date + "\n"
		stringToTweet += "With data from the GR Gov API\n"
		stringToTweet += "#koronoios #covid19GR #CovidGR #emvolio #COVID19 #CoronavirusVaccine"
		fmt.Println(stringToTweet)
	} else {
		stringToTweet += "Εώς " + date + "\n"
		stringToTweet += "Με δεδομένα από το GR Gov API \n"
		stringToTweet += "#emvolio #covid #COVID19gr #κορονοϊός #εμβόλιο #εμβολιασμος #κορωνοιος"
		fmt.Println(stringToTweet)
	}
	sendTweet(stringToTweet)
}

func sendTweet(stringToTweet string) {
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	fmt.Printf("%+v\n", creds)

	client, err := getClient(&creds)
	if err != nil {
		fmt.Println("Error getting Twitter Client")
		fmt.Println(err)
	}
	tweet, resp, err := client.Statuses.Update(stringToTweet, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", tweet)
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	fmt.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

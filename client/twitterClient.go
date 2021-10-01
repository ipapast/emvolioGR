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
	solidBarsToPrint := math.Round(dataValue / 6.67)
	emptyBarsToPrint := 15 - int(solidBarsToPrint)

	dataToAdd += textValue
	for solidBarsToPrint > 0 {
		dataToAdd += "▓"
		solidBarsToPrint -= 1
	}
	for emptyBarsToPrint > 0 {
		dataToAdd += "░"
		emptyBarsToPrint -= 1
	}
	dataToAdd += " " + fmt.Sprint(dataValue) + "%\n\n"
	return dataToAdd
}

func SourceAndSendTweet(stringToTweet string, language string) {
	// Github Action is getting triggered in the morning, when we have data available until yesterday
	yesterday := time.Now().Add(-24 * time.Hour)
	yesterdayDate := fmt.Sprint(yesterday.Format("02/01/2006"))

	if language == "en" {
		stringToTweet += "As of " + yesterdayDate + "\n"
		stringToTweet += "With data from the GR Gov API\n"
		stringToTweet += "#koronoios #covid19GR #CovidGR #emvolio #COVID19 #CoronavirusVaccine"
		fmt.Println(stringToTweet)
	} else {
		stringToTweet += "Εώς " + yesterdayDate + "\n"
		stringToTweet += "Με δεδομένα από το GR Gov API \n"
		stringToTweet += "#emvolio #covid #COVID19gr #κορονοϊός #εμβόλιο #εμβολιασμος #κορωνοιος"
		fmt.Println(stringToTweet)
	}
	sendTweet(stringToTweet)
}

func sendTweet(stringToTweet string) {
	credentials := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	fmt.Printf("%+v\n", credentials)

	client, err := getClient(&credentials)
	if err != nil {
		fmt.Println("Error getting Twitter Client")
		os.Exit(1)
	}
	tweet, resp, err := client.Statuses.Update(stringToTweet, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", tweet)
}

func getClient(credentials *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(credentials.ConsumerKey, credentials.ConsumerSecret)
	token := oauth1.NewToken(credentials.AccessToken, credentials.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		fmt.Println("Could not verify Twitter credentials")
		return nil, err
	}

	fmt.Printf("User's ACCOUNT:\n%+v\n", user)
	return client, nil
}

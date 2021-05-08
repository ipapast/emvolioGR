package main

import (
	// "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
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

func getData() {
	fmt.Println("Getting data...")

	url := "https://data.gov.gr/api/v1/query/mdg_emvolio"
	gov_token := os.Getenv("GOV_DATA_TOKEN")

	token := fmt.Sprint("Token ", gov_token)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading the response bytes:", err)
	}
	// fmt.Println(string([]byte(body)))
	data := []byte(body)

	jsonFile, _ := os.Create("vaccinations_regions.json")
	defer jsonFile.Close()

	_, err2 := jsonFile.Write(data)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Done writing to file")

	type DataPerArea struct {
		ReferenceDate        string `json:"referencedate"`
		TotalVaccinations    int    `json:"totalvaccinations"`
		TotalDistinctPersons int    `json:"totaldistinctpersons"`
	}

	type AllData []DataPerArea
	var allData AllData

	err = json.Unmarshal([]byte(body), &allData)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println(allData)

	for i := range allData {
		fmt.Println(allData[i].ReferenceDate)

	}

	// map = { "01/05/2021": {totalpeople: "50", totalVaccs: "120" }

	// dataMap := make(map[string]struct{})
	// for k, _ := range allData {
	// 	dataMap["k"]
	// 	fmt.Printf("key[%s] \n", k)
	// }

}

func main() {
	getData()

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
		tweet, resp, err := client.Statuses.Update("A Test Tweet", nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", resp)
		fmt.Printf("%+v\n", tweet)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
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

	total_vacs := make(map[time.Time]int)
	total_people_vac := make(map[time.Time]int)
	total_people_vac_fully := make(map[time.Time]int)
	population_of_gr := 8868536
	var percentage_1st_dose []float64
	var percentage_2nd_dose []float64

	for _, elem := range allData {
		date := elem.ReferenceDate[0:10]
		parsed_date, _ := time.Parse("2006-01-02", date)
		// fmt.Println(parsed_date)
		total_vacs[parsed_date] += elem.TotalVaccinations
		total_people_vac[parsed_date] += elem.TotalDistinctPersons
		total_people_vac_fully[parsed_date] = total_vacs[parsed_date] - total_people_vac[parsed_date]

		perc := (float64(total_people_vac[parsed_date]) / float64(population_of_gr)) * 100
		percentage_1st_dose = append(percentage_1st_dose, math.Round(perc*100)/100)
		perc2 := (float64(total_people_vac_fully[parsed_date]) / float64(population_of_gr)) * 100
		percentage_2nd_dose = append(percentage_2nd_dose, math.Round(perc2*100)/100)
	}

	sort.Slice(percentage_1st_dose, func(i, j int) bool { return percentage_1st_dose[i] < percentage_1st_dose[j] })
	sort.Slice(percentage_2nd_dose, func(i, j int) bool { return percentage_2nd_dose[i] < percentage_2nd_dose[j] })
	// get last item
	fmt.Println(percentage_1st_dose[len(percentage_1st_dose)-1])
	fmt.Println(percentage_2nd_dose[len(percentage_2nd_dose)-1])
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

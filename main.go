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

	"github.com/ipapast/emvolioGR/client"
)

func getCovidVaccinations() []byte {
	url := "https://data.gov.gr/api/v1/query/mdg_emvolio"
	govToken := os.Getenv("GOV_DATA_TOKEN")

	token := fmt.Sprint("Token ", govToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error on getting data from Gov GR.\n[ERROR] -", err)
	}

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
	// writeToFile(body)

	return body
}

func writeToFile(body []byte) {
	jsonFile, _ := os.Create("../data/vaccinations_regions.json")
	defer jsonFile.Close()

	_, err2 := jsonFile.Write(body)
	if err2 != nil {
		fmt.Println(err2)
	}
}
func transformData(res []byte) {
	type DataPerArea struct {
		ReferenceDate        string `json:"referencedate"`
		TotalVaccinations    int    `json:"totalvaccinations"`
		TotalDistinctPersons int    `json:"totaldistinctpersons"`
	}

	type AllData []DataPerArea
	var allData AllData

	err := json.Unmarshal((res), &allData)
	if err != nil {
		fmt.Println(err)
	}

	totalVacs := make(map[time.Time]int)
	totalPeopleVac := make(map[time.Time]int)
	totalPeopleVacFully := make(map[time.Time]int)
	populationOfGr := 8868536
	var percentage1stDose []float64
	var percentage2ndDose []float64

	for _, elem := range allData {
		date := elem.ReferenceDate[0:10]
		parsedDate, _ := time.Parse("2006-01-02", date)
		totalVacs[parsedDate] += elem.TotalVaccinations
		totalPeopleVac[parsedDate] += elem.TotalDistinctPersons
		totalPeopleVacFully[parsedDate] = totalVacs[parsedDate] - totalPeopleVac[parsedDate]

		percentage1st := (float64(totalPeopleVac[parsedDate]) / float64(populationOfGr)) * 100
		percentage1stDose = append(percentage1stDose, math.Round(percentage1st*100)/100)
		percantage2nd := (float64(totalPeopleVacFully[parsedDate]) / float64(populationOfGr)) * 100
		percentage2ndDose = append(percentage2ndDose, math.Round(percantage2nd*100)/100)
	}

	sort.Slice(percentage1stDose, func(i, j int) bool { return percentage1stDose[i] < percentage1stDose[j] })
	sort.Slice(percentage2ndDose, func(i, j int) bool { return percentage2ndDose[i] < percentage2ndDose[j] })
	// get last item
	latest1stDose := percentage1stDose[len(percentage1stDose)-1]
	latest2ndDose := percentage2ndDose[len(percentage2ndDose)-1]

	stringToTweet := ""
	stringToTweet += client.AddDataToTweet(latest1stDose, "1st dose of vaccine progress in Greece: \n\n")
	stringToTweet += client.AddDataToTweet(latest2ndDose, "2nd dose of vaccine progress in Greece: \n\n")
	stringToTweetGR := ""
	stringToTweetGR += client.AddDataToTweet(latest1stDose, "Ποσοστό ατόμων με 1η δόση εμβολίου: \n\n")
	stringToTweetGR += client.AddDataToTweet(latest2ndDose, "Ποσοστό ατόμων με 2η δόση εμβολίου: \n\n")
	client.SourceAndSendTweet(stringToTweet, "en")
	client.SourceAndSendTweet(stringToTweetGR, "gr")
}

func main() {
	res := getCovidVaccinations()
	transformData(res)
}

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

)

func getData() {
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
	latest_1st_dose := (percentage_1st_dose[len(percentage_1st_dose)-1])
	latest_2nd_dose := (percentage_2nd_dose[len(percentage_2nd_dose)-1])

	stringToTweet := ""
	stringToTweet += AddDataToTweet(latest_1st_dose, "1st dose of vaccine progress in Greece: \n\n")
	stringToTweet += AddDataToTweet(latest_2nd_dose, "2nd dose of vaccine progress in Greece: \n\n")
	stringToTweetGR := ""
	stringToTweetGR += AddDataToTweet(latest_1st_dose, "Ποσοστό ατόμων με 1η δόση εμβολίου: \n\n")
	stringToTweetGR += AddDataToTweet(latest_2nd_dose, "Ποσοστό ατόμων με 2η δόση εμβολίου: \n\n")
	SourceAndSendTweet(stringToTweet, "en")
	SourceAndSendTweet(stringToTweetGR, "gr")
}

func main() {
	getData()
}

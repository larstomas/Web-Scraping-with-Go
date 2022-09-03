package main

// https://www.scrapingbee.com/blog/web-scraping-go/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"

	"github.com/gocolly/colly"
)

type DataRow struct {
	Name  string `json:"name"`
	Candy string `json:"candy"`
	Eaten int    `json:"eaten"`
}

var theData []DataRow
var top3 []DataRow

func WebScraper() {

	c := colly.NewCollector()

	// Selector: #top\\.customers > tbody > tr:nth-child(1) > td:nth-child(1)
	c.OnHTML("table#top\\.customers > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			count, _ := strconv.Atoi(el.ChildText("td:nth-child(3)"))
			dataRow := DataRow{
				Name:  el.ChildText("td:nth-child(1)"),
				Candy: el.ChildText("td:nth-child(2)"),
				Eaten: count,
			}
			theData = append(theData, dataRow)
		})
		fmt.Println("Scrapping Complete")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://candystore.zimpler.net/")
}

func writeJSON(theData []DataRow) {
	json_file, err := json.MarshalIndent(theData, "", " ")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = ioutil.WriteFile("TheData.json", json_file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func makeTop3(theData []DataRow) {
	// Sort data
	sort.Slice(theData[:], func(i, j int) bool {
		return theData[i].Eaten > theData[j].Eaten
	})

	for i := 0; i < 3; i++ {
		top3 = append(top3, theData[i])
	}
}

func main() {
	WebScraper()
	makeTop3(theData)
	fmt.Printf("%+v\n", top3)
	writeJSON(top3)
}

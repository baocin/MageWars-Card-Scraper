package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Card holds all info about a MageWars card
type Card struct {
	Name, CardType, Set, Levels, ImageURL, CardCode, Cost, RevealCost string
	SubTypes, Schools                                                 []string
}

func main() {
	//The website with the table of cards
	url := "http://forum.arcanewonders.com/sbb/database.php"
	//

	// bytes, response := downloadRaw(url)
	// html := string(bytes)

	cards := getCards(url)
	for _, card := range cards {
		fmt.Println(card)
		downloadFile(card.ImageURL, card.Name + " " + card.CardCode + ".jpg")

		// Save the JSON of this card to filename
		cardJSON, _ := json.Marshal(card)
		fmt.Println(string(cardJSON))
		ioutil.WriteFile(card.Name + " " + card.CardCode + ".json", cardJSON, os.ModePerm)
	}

	//Save all cards into single json
	cardsJSON, _ := json.Marshal(cards)
	fmt.Println(string(cardsJSON))
	ioutil.WriteFile("AllCards.json", cardsJSON, os.ModePerm)

}

func getCards(url string) []Card {
	baseImageURL := "http://forum.arcanewonders.com/cards/"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	var cards []Card

	// Find the cards in the html
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		cardCode, isValidCard := s.Attr("data-code")
		if isValidCard {
			set, _ := s.Attr("data-set")

			columns := s.Children()
			newCard := Card{
				Name:       columns.First().Text(),
				CardType:   columns.Eq(1).Text(),
				SubTypes:   strings.Split(columns.Eq(2).Text(), ", "),
				Schools:    strings.Split(columns.Eq(3).Text(), ", "),
				Levels:     columns.Eq(4).Text(),
				Cost:       columns.Eq(5).Text(),
				RevealCost: columns.Eq(6).Text(),
				Set:        set,
				ImageURL:   baseImageURL + cardCode + ".jpg",
				CardCode:   cardCode,
			}

			cards = append(cards, newCard)
		}
	})
	return cards
}

func downloadRaw(url string) (bytes []byte, response *http.Response) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
		log.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()

	byteResult, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bytes = byteResult
	response = resp

	return
}

func downloadFile(url string, filename string) {
	bytes, _ := downloadRaw(url)
	ioutil.WriteFile(filename, bytes, os.ModePerm)
}

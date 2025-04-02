package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type card struct {
	Name      string `json:"name"`
	Cost      string `json:"cost"`
	Power     string `json:"power"`
	Counter   string `json:"counter"`
	Color     string `json:"color"`
	Type      string `json:"type"`
	Effect    string `json:"effect"`
	CardSet   string `json:"set"`
	Attribute string `json:"attribute"`
	CardNo    int    `json:"cardNo"`
	ImgPath   string `json:"imgPath"`
	Info      string `json:"info"`
}

type store struct {
	data []card
}

var dataStore = &store{}

func loadData() error {
	file, err := os.Open("../data.json")
	if err != nil {
		return fmt.Errorf("failed to open data.json: %w", err)
	}
	defer file.Close()

	var cards []card
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cards); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	dataStore.data = cards
	return nil
}

func getByName(c *gin.Context) {
	name := c.Param("name")
	var results []card

	for _, card := range dataStore.data {
		if strings.Contains(strings.ToLower(card.Name), strings.ToLower(name)) {
			results = append(results, card)
		}
	}

	c.IndentedJSON(http.StatusOK, results)
	writeLog(c.ClientIP())
}

func getByID(c *gin.Context) {
	var results []card

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	for _, card := range dataStore.data {
		if id == card.CardNo {
			results = append(results, card)
		}
	}

	c.IndentedJSON(http.StatusOK, results)
	writeLog(c.ClientIP())
}

func getBySet(c *gin.Context) {
	var results []card
	set := c.Param("set")

	for _, card := range dataStore.data {
		if strings.Contains(strings.ToLower(card.CardSet), strings.ToLower(set)) {
			results = append(results, card)
		}
	}

	c.IndentedJSON(http.StatusOK, results)
	writeLog(c.ClientIP())
}

func getByColor(c *gin.Context) {
	var results []card
	var colors []string = strings.Split(strings.ToLower(c.Param("color")), "-")

	colorSet := make(map[string]struct{}, len(colors))
	for _, color := range colors {
		colorSet[color] = struct{}{}
	}

	for _, card := range dataStore.data {
		cardColors := strings.Split(strings.ToLower(card.Color), "/")

		cardColorSet := make(map[string]struct{}, len(cardColors))
		for _, c := range cardColors {
			cardColorSet[c] = struct{}{}
		}

		for color := range colorSet {
			if _, found := cardColorSet[color]; found {
				results = append(results, card)
				break
			}
		}
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getByColorSqL(c *gin.Context) {
	dbCon := connectToDB()
	defer dbCon.Close()

	var colors []string = strings.Split(strings.ToLower(c.Param("color")), "-")

	results, err := queryByColor(colors, dbCon)
	if err != nil {
		fmt.Println(err)
	}

	c.IndentedJSON(http.StatusOK, results)
}

// Consider adding pagnation and limit params for handling
func getAllCards(c *gin.Context) {
	c.JSON(http.StatusOK, dataStore.data)
	writeLog(c.ClientIP())
}

// Consider a check for out of bounds and handling
func getAllCardsPag(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit")) //amount of cards to fetch
	if err != nil {
		log.Fatal(err)
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.Fatal(err)
	}

	var results []card

	lastIndex := limit * page

	for i := lastIndex - limit; i < lastIndex && len(results) < limit; i++ {
		results = append(results, dataStore.data[i])
	}

	c.IndentedJSON(http.StatusOK, results)
}

func main() {
	err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Static("/assets", "../assets")

	router.GET("/all", getAllCards)
	router.GET("/all/:limit/:page", getAllCardsPag)
	router.GET("/name/:name", getByName)
	router.GET("/id/:id", getByID)
	router.GET("/set/:set", getBySet)
	router.GET("/color/:color", getByColor)

	router.GET("/colorSQL/:color", getByColorSqL)

	//Unhandled err
	router.Run("localhost:8080")
}

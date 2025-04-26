package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

type databaseConnection struct {
	db *sql.DB
}

var dbConn = databaseConnection{}
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
}

func getByID(c *gin.Context) {
	var results []card

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error with id param": err.Error()})
	}

	for _, card := range dataStore.data {
		if id == card.CardNo {
			results = append(results, card)
		}
	}

	c.IndentedJSON(http.StatusOK, results)
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
}

func getBySetSQL(c *gin.Context) {
	var results []card
	set := c.Param("set")

	results, err := queryBySet(set, dbConn.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error with query": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getByColor(c *gin.Context) {
	var results []card
	var colors = strings.Split(strings.ToLower(c.Param("color")), "-")

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
	var colors = strings.Split(strings.ToLower(c.Param("color")), "-")

	results, err := queryByColor(colors, dbConn.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error with query": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getByIdSQL(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error with id param": err.Error()})
		return
	}

	results, err := queryByID(id, dbConn.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error during query": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func getAllCards(c *gin.Context) {
	c.JSON(http.StatusOK, dataStore.data)
}

func getAllCardsPag(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: issue with limit param": err.Error()})
		return
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: issue with page param": err.Error()})
		return
	}

	var results []card
	lastIndex := limit * page

	for i := lastIndex - limit; i < lastIndex && len(results) < limit; i++ {
		results = append(results, dataStore.data[i])
	}

	c.IndentedJSON(http.StatusOK, results)
}

// TODO: Testing needed
func getAllCardsPageSQL(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: issue with limit param": err.Error()})
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: issue with page param": err.Error()})
	}

	results, err := queryPag(limit, page, dbConn.db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error during query": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, results)
}

// TODO: Add pooling here to handling lots of connections as once
func main() {
	err := loadData()
	if err != nil {
		panic(err)
	}

	dbConn.db = connectToDB()
	defer dbConn.db.Close()

	router := gin.Default()

	router.Static("/assets", "../assets")

	//JSON Parsing endpoints
	router.GET("/all", getAllCards)
	router.GET("/all/:limit/:page", getAllCardsPag)
	router.GET("/name/:name", getByName)
	router.GET("/id/:id", getByID)
	router.GET("/set/:set", getBySet)
	router.GET("/color/:color", getByColor)

	//MySQL endpoints
	router.GET("/colorSQL/:color", getByColorSqL)
	router.GET("/setSQL/:set", getBySetSQL)
	router.GET("/idSQL/:id", getByIdSQL)
	router.GET("/allSQL/:limit/:page", getAllCardsPageSQL)

	err = router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectToDB() *sql.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8mb4", user, password, host, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

// TODO: Testing needed
func parseResponse(rows *sql.Rows) ([]card, error) {
	var cards []card
	for rows.Next() {
		var cardItem card
		if err := rows.Scan(
			&cardItem.CardNo,
			&cardItem.Name,
			&cardItem.Cost,
			&cardItem.Power,
			&cardItem.Counter,
			&cardItem.Color,
			&cardItem.Type,
			&cardItem.Effect,
			&cardItem.CardSet,
			&cardItem.Attribute,
			&cardItem.ImgPath,
			&cardItem.Info,
		); err != nil {
			return nil, err
		}

		cards = append(cards, cardItem)
	}

	return cards, nil
}

func queryByColor(colors []string, db *sql.DB) ([]card, error) {
	query := "Select distinct * from cards where "

	conditions := make([]string, len(colors))
	args := make([]interface{}, len(colors))

	for i, color := range colors {
		conditions[i] = "color LIKE ?"
		args[i] = "%" + color + "%"
	}

	query += strings.Join(conditions, " OR ")

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cards, err := parseResponse(rows)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func queryBySet(set string, db *sql.DB) ([]card, error) {
	query := "Select * from cards where CardSet LIKE ?"
	arg := "%" + set + "%"

	rows, err := db.Query(query, arg)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cards, err := parseResponse(rows)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func queryByID(id int, db *sql.DB) ([]card, error) {
	query := "Select * from cards where CardNo = ?"

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cards, err := parseResponse(rows)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

// FIXME: This is off by 1.
func queryPag(limit int, page int, db *sql.DB) ([]card, error) {
	query := "select * from cards limit ? offset ?"

	rows, err := db.Query(query, limit, page)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cards, err := parseResponse(rows)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return cards, nil
}

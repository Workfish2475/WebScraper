package main

import (
	"database/sql"
	"time"
)

func connectToDB(dbHost string) *sql.DB {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func queryByColor(color string, db *sql.DB) ([]card, error) {
	var cards []card

	query := "Select * from cards where color like ?"
	searchTerm := "%" + color + "%"

	rows, err := db.Query(query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cardItem card
		if err := rows.Scan(&cardItem.Name, &cardItem.Cost, &cardItem.Power, &cardItem.Counter, &cardItem.Color, &cardItem.Type, &cardItem.Effect, &cardItem.Set, &cardItem.Attribute, &cardItem.CardNo, &cardItem.Info); err != nil {
			return nil, err
		}

		cards = append(cards, cardItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

package main

import (
	"database/sql"
)

type logItem struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

func logInteraction(db *sql.DB) {

}

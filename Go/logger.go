package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type logItem struct {
	Time string `json:"time"`
	IP   string `json:"ip"`
}

func writeLog(ip string) {
	var logs []logItem

	logFile, err := os.Open("logs.json")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	decoder := json.NewDecoder(logFile)
	if err := decoder.Decode(&logs); err != nil {
		fmt.Errorf("An error occurred: ", err)
	}

	logs = append(logs, logItem{Time: time.Now().Format(time.RFC3339), IP: ip})

	updateData, err := json.MarshalIndent(logs, "", "\t")
	if err != nil {
		fmt.Errorf("An error occurred: ", err)
	}

	err = os.WriteFile("logs.json", updateData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

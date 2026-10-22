package main

import (
	"log"
	"time"
	"user-service/cmd"
)

func main() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}

	time.Local = loc

	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed to execute: %v", err)
	}
}

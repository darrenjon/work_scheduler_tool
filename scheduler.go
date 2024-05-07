package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func callAPI() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiUrl := os.Getenv("EKG_SCAN_API")
	// Get current date and format it to "YYYYMD"
	now := time.Now()
	dateStrToday := fmt.Sprintf("%d%d%d", now.Year(), int(now.Month()), now.Day())

	// Define the request payload
	payload := map[string]interface{}{
		"date_str": []string{dateStrToday},
		"backfill": true,
		"last":     true,
		"test":     false,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error encoding payload: %v", err)
	}

	// Print the request payload
	log.Printf("API Request Payload: %s", string(payloadBytes))

	// Create a new request using http
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Log the response status
	log.Printf("API Response Status: %s", resp.Status)

	// Decode the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	// Print the response
	log.Printf("API Response: %v", result)
}

func scheduleDaily(f func()) {
	now := time.Now()
	// Calculate next occurrence of the specified time
	next := now.Truncate(time.Hour).Add(time.Hour)
	if next.Before(now) {
		next = next.Add(2 * time.Hour)
	}
	time.Sleep(next.Sub(now)) // Sleep until the next occurrence
	go f()

	for {
		// Sleep for 2 hour before scheduling the next occurrence
		time.Sleep(2 * time.Hour)
		go f()
	}
}

func main() {
	log.Println("Starting scheduler...")
	scheduleDaily(callAPI)
}

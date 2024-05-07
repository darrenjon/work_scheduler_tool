package main

import (
	"bytes"
	"encoding/json"
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
	dateStrToday := now.Format("2006012")

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
}

func scheduleDaily(f func()) {
	now := time.Now()
	// Calculate next occurrence of the specified time
	next := now.Truncate(time.Hour).Add(time.Hour)
	if next.Before(now) {
		next = next.Add(2 * time.Hour)
	}
	time.Sleep(next.Sub(now)) // Sleep until the next occurrence
	go f()                    // Run the function in a new goroutine

	for {
		// Sleep for 24 hours, then run again
		time.Sleep(2 * time.Hour)
		go f()
	}
}

func main() {
	// Schedule the API call function daily at 7:00 AM
	scheduleDaily(callAPI)
}

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

func callAPI(apiUrl string, payloadBytes []byte) {
	// Create a new request using http
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("API Response Status: %s", resp.Status)

	// Decode the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return
	}

	// Print the response
	log.Printf("API Response: %v", result)
}

func scheduleEkgScanApi() {
	log.Println("Running EKG_SCAN_API...")
	ekgScanApiUrl := os.Getenv("EKG_SCAN_API")

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
		log.Printf("Error encoding payload: %v", err)
	}

	log.Printf("API Request Payload: %s", string(payloadBytes))

	// Call EKG_SCAN_API
	callAPI(ekgScanApiUrl, payloadBytes)
}

func scheduleIpdHandoverApi() {
	log.Println("Running HANDOVER_UPDATE_API...")
	ipdHandoverApiUrl := os.Getenv("HANDOVER_UPDATE_API")

	callAPI(ipdHandoverApiUrl, nil)
}

func scheduleIpdHandoverApiShh() {
	log.Println("Running HANDOVER_UPDATE_API_SHH...")
	ipdHandoverApiUrlShh := os.Getenv("HANDOVER_UPDATE_API_SHH")

	callAPI(ipdHandoverApiUrlShh, nil)
}

func scheduleDaily(f func(), interval time.Duration, start time.Duration) {
	now := time.Now()
	next := now.Truncate(time.Hour).Add(start)
	if next.Before(now) {
		next = next.Add(interval)
	}

	time.Sleep(next.Sub(now))
	go f()

	for {
		time.Sleep(interval)
		go f()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Starting scheduler...")

	go scheduleDaily(scheduleEkgScanApi, 1*time.Hour, 0)
	go scheduleDaily(scheduleIpdHandoverApi, 30*time.Minute, 30*time.Minute)
	go scheduleDaily(scheduleIpdHandoverApiShh, 30*time.Minute, 30*time.Minute)

	// Wait indefinitely
	select {}
}

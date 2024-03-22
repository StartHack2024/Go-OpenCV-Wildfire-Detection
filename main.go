package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
)

type FireDetection struct {
	DateTime  string  `json:"datetime"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	http.HandleFunc("/fireAlert", sendEmailHandler)
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Failed to close request body: %v", err)
		}
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	var fireDetection FireDetection
	if err := json.Unmarshal(body, &fireDetection); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	if fireDetection.DateTime == "" || fireDetection.Latitude == 0 || fireDetection.Longitude == 0 {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	emailBody := fmt.Sprintf("Forest fire detected!\n\nDateTime: %s\nLatitude: %f\nLongitude: %f",
		fireDetection.DateTime, fireDetection.Latitude, fireDetection.Longitude)

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "sir.alexander.dark@gmail.com" // Use environment variables or config
	smtpPassword := "iblr flfv cbad iyds"          // Use environment variables or config
	recipient := "sir_alexiner@hotmail.com"

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Forest Fire Alert!\r\n\r\n%s", smtpUsername, recipient, emailBody)
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	if err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, smtpUsername, []string{recipient}, []byte(message)); err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Email sent successfully")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

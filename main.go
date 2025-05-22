package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// Shared variable and mutex to protect concurrent access
var (
	currentMoisture int
	mu              sync.RWMutex
)

func main() {
	// Starts serial reading in the background
	go readSerialWindows("COM3") // <-- Replace COM3 with your actual port, this is however different on Linux and MacOS

	// Sets up web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		data := currentMoisture
		mu.RUnlock()

		// Parses and render the HTML template
		tmpl, err := template.ParseFiles("template.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// readSerialWindows reads serial data from Arduino using tarm/serial (cross-platform)
func readSerialWindows(portName string) {
	// Configure the serial port
	config := &serial.Config{
		Name:        portName, // COM3, COM4, etc.
		Baud:        9600,     // Match your Arduino's baud rate, the default is 9600
		ReadTimeout: time.Second * 5,
	}

	// Open the serial port
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer s.Close()

	// Use a scanner to read line-by-line
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		line := scanner.Text()
		if value, err := strconv.Atoi(line); err == nil {
			mu.Lock()
			currentMoisture = value
			mu.Unlock()
		}
	}
}

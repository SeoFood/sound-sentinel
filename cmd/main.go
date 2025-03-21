package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"sound-sentinel/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Sound Sentinel...")

	// Lade .env Datei
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	detector, err := app.NewSoundDetector()
	if err != nil {
		log.Fatalf("Error initializing Sound Detector: %v", err)
	}

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	log.Printf("Audio source: %s", detector.AudioSource)
	log.Printf("Sound threshold: %.2f", detector.SoundThreshold)
	log.Printf("MQTT broker: %s", detector.MQTTBroker)

	err = detector.Run()
	if err != nil {
		log.Fatalf("Error running Sound Detector: %v", err)
	}
}

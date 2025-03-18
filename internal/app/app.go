package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

// SoundDetector represents the main application
type SoundDetector struct {
	AudioSource      string
	SoundThreshold   float64
	DetectionTimeout int
	MQTTBroker       string
	MQTTTopic        string
	MQTTUsername     string
	MQTTPassword     string
	MQTTClientID     string
	client           mqtt.Client
	soundDetected    bool
	lastDetection    time.Time
}

// NewSoundDetector creates a new instance of SoundDetector
func NewSoundDetector() (*SoundDetector, error) {
	// Lade .env Datei (ignoriere Fehler, wenn die Datei nicht existiert)
	_ = godotenv.Load()

	threshold, err := strconv.ParseFloat(getEnv("SOUND_THRESHOLD", "0.03"), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid SOUND_THRESHOLD: %v", err)
	}

	timeout, err := strconv.Atoi(getEnv("DETECTION_TIMEOUT", "5"))
	if err != nil {
		return nil, fmt.Errorf("Invalid DETECTION_TIMEOUT: %v", err)
	}

	return &SoundDetector{
		AudioSource:      getEnv("AUDIO_SOURCE", "default"),
		SoundThreshold:   threshold,
		DetectionTimeout: timeout,
		MQTTBroker:       getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTTopic:        getEnv("MQTT_TOPIC", "homeassistant/sensor/sound_sentinel/state"),
		MQTTUsername:     getEnv("MQTT_USERNAME", ""),
		MQTTPassword:     getEnv("MQTT_PASSWORD", ""),
		MQTTClientID:     getEnv("MQTT_CLIENT_ID", "sound-sentinel"),
		soundDetected:    false,
		lastDetection:    time.Time{},
	}, nil
}

// Run starts the sound detection
func (sd *SoundDetector) Run() error {
	err := sd.connectMQTT()
	if err != nil {
		return err
	}

	err = sd.monitorSound()
	return err
}

// connectMQTT establishes a connection to the MQTT broker
func (sd *SoundDetector) connectMQTT() error {
	opts := mqtt.NewClientOptions().
		AddBroker(sd.MQTTBroker).
		SetClientID(sd.MQTTClientID).
		SetAutoReconnect(true).
		SetConnectionLostHandler(sd.onConnectionLost).
		SetOnConnectHandler(sd.onConnect)

	if sd.MQTTUsername != "" {
		opts.SetUsername(sd.MQTTUsername)
		opts.SetPassword(sd.MQTTPassword)
	}

	sd.client = mqtt.NewClient(opts)
	token := sd.client.Connect()
	token.Wait()
	return token.Error()
}

// onConnect is the handler for MQTT connections
func (sd *SoundDetector) onConnect(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	sd.publishStatus("OFF")
}

// onConnectionLost is the handler for lost connections
func (sd *SoundDetector) onConnectionLost(client mqtt.Client, err error) {
	log.Printf("MQTT connection lost: %v", err)
}

// publishStatus sends the current status to MQTT
func (sd *SoundDetector) publishStatus(status string) {
	token := sd.client.Publish(sd.MQTTTopic, 0, true, status)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Error publishing: %v", token.Error())
	}
}

// monitorSound starts FFmpeg and monitors the audio output
func (sd *SoundDetector) monitorSound() error {
	for {
		cmd := sd.createFFmpegCommand()

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Error setting up stdout pipe: %v", err)
		}

		err = cmd.Start()
		if err != nil {
			log.Printf("Error starting FFmpeg: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			sd.processAudioData(line)
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("FFmpeg exited with error: %v", err)
		}

		time.Sleep(2 * time.Second)
		log.Println("Restarting FFmpeg...")
	}
}

// createFFmpegCommand creates the FFmpeg command
func (sd *SoundDetector) createFFmpegCommand() *exec.Cmd {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-i", sd.AudioSource,
		"-af", "volumedetect",
		"-f", "null",
		"-",
	}
	return exec.Command("ffmpeg", args...)
}

// processAudioData processes audio output from FFmpeg
func (sd *SoundDetector) processAudioData(line string) {
	if !strings.Contains(line, "max_volume") {
		return
	}

	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return
	}

	volumeStr := strings.TrimSpace(parts[1])
	volumeStr = strings.Replace(volumeStr, " dB", "", 1)

	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		log.Printf("Error parsing volume: %v", err)
		return
	}

	// Convert negative volume values to absolute values
	volumeLevel := -volume

	log.Printf("Detected volume: %.2f dB (threshold: %.2f)", volumeLevel, sd.SoundThreshold)

	// Sound detection based on threshold
	if volumeLevel >= sd.SoundThreshold {
		sd.lastDetection = time.Now()
		if !sd.soundDetected {
			sd.soundDetected = true
			sd.publishStatus("ON")
			log.Println("Sound detected: ON")
		}
	} else if sd.soundDetected && time.Since(sd.lastDetection).Seconds() > float64(sd.DetectionTimeout) {
		sd.soundDetected = false
		sd.publishStatus("OFF")
		log.Println("No sound detection: OFF")
	}
}

// getEnv gets environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

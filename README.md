# Sound Detection - FFmpeg Audio Sound Detection to MQTT

This project provides a Docker image that uses FFmpeg to detect sound from audio sources and sends MQTT messages to Home Assistant when sound is detected.

## Features

- Continuous monitoring and detection of sound from various audio sources (microphone, streams, files)
- Real-time sound level measurement and reporting
- Configurable thresholds for sound detection
- Immediate notification via MQTT when audio conditions change
- Simple configuration via environment variables
- Robust error handling with automatic recovery and reconnection
- Comprehensive logging for troubleshooting

## Technologies

- Golang 1.18+
- FFmpeg for audio analysis
- MQTT for messaging
- Docker for containerization

## Installation

### Prerequisites

- Docker and Docker Compose
- MQTT Broker (e.g., Mosquitto in Home Assistant)

### Setup

## Configuration

All configuration options are controlled via environment variables in the `.env` file:

## Troubleshooting

If the sound detection is not working as expected:

1. Check the logs for any errors:

   ```bash
   docker-compose logs -f
   ```

2. Verify your audio source is working correctly
3. Try adjusting the sound threshold based on your environment noise level

4. Make sure your MQTT connection is working properly

### Testing Your Audio Source

You can use the included test script to check if your audio source is correctly detected and to determine the appropriate sound threshold:

1. Make sure you have set up your `.env` file with the correct `AUDIO_SOURCE`.

2. Run the test command:

   ```bash
   script/test.sh
   ```

3. If using a Docker container, you may need to run the test inside the container:

   ```bash
   docker exec -it sound-sentinel /bin/bash
   ```

## Project Structure

- **cmd/**: Contains the entry point of the application.

  - **main.go**: The main function that starts the application.

- **internal/**: Contains the core logic of the application.

  - **app/**: Contains the main application structure and methods.
    - **app.go**: Defines the application logic.

- **pkg/**: Contains utility functions that can be used throughout the application.

  - **utils/**: Contains helper functions.
    - **helpers.go**: Exports general-purpose functions.

- **.vscode/**: Contains configuration files for the development environment.

  - **settings.json**: Configuration settings for Go projects.
  - **extensions.json**: Recommended extensions for enhancing the development experience.

- **go.mod**: Module definition for the Go project, listing dependencies and Go version.

- **go.sum**: Contains checksums for dependencies to ensure the correct versions are used.

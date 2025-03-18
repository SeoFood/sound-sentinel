#!/bin/bash

source .env 2>/dev/null || echo "Warning: No .env file found"

echo "Testing audio analysis..."
echo "Audio source: ${AUDIO_SOURCE:-default}"

ffmpeg -hide_banner -loglevel info -t 10 -i "${AUDIO_SOURCE:-default}" -af volumedetect -f null - 2>&1 | grep max_volume

echo "Test completed. Note the 'max_volume' values to set an appropriate threshold."
echo "A higher (less negative) dB value means louder sound."
echo "Set SOUND_THRESHOLD in your .env file based on these values."

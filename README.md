# midi2homeassistant

This application lets you control smart lights via a MIDI controller using Home Assistant.

## Overview

The application listens for MIDI events from your controller and sends corresponding commands to Home Assistant to control your smart lights.

## Setup Instructions

### 1. Locate Your MIDI Device

On Linux, the USB MIDI device buffer is typically found under `/dev/snd/`. For example, you might find your device at `/dev/snd/midiC2D0`. Use a command like:

```bash
ls /dev/snd/
cat /dev/snd/midiC2D0  # If this is the correct buffer, some text will appear whenever you press a key on your MIDI device.
```

# Home Assistant configuration

To configure the application, update the following environment variables in your .env file:
```ini
IP="your-home-assistant-ip"
PORT=8123
ENDPOINT="/api"
TOKEN="your_long_lived_access_token"
ENTITY_ID="your_light_entity_id"
MIDI_PATH="/dev/snd/your_midi_device"
BRIGHTNESS_INCREMENT=50
```

## Adding your own configurations

You can add your own configurations by reading the Home Assistant documentation and observing/experimenting with the MIDI buffer.

# I made a [video](https://www.youtube.com/watch?v=0_4cr9nkL7k&ab_channel=HeyDayTime) about this project
![image](https://github.com/user-attachments/assets/76ba490b-2c66-4c5d-bc31-ee867a3a1158)
> Yeah, I’ve got a buzzcut in the video. Don’t ask why.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	IP        string
	PORT      string
	ENDPOINT  string
	TOKEN     string
	ENTITY_ID string
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	myAuth := Auth{
		IP:        os.Getenv("IP"),
		PORT:      os.Getenv("PORT"),
		ENDPOINT:  os.Getenv("ENDPOINT"),
		TOKEN:     os.Getenv("TOKEN"),
		ENTITY_ID: os.Getenv("ENTITY_ID"),
	}

	BRIGHTNESS_INCREMENT, err := strconv.Atoi((os.Getenv("BRIGHTNESS_INCREMENT")))
	if err != nil {
		log.Fatal(err)
	}

	devicePath := os.Getenv("MIDI_PATH")

	file, err := os.Open(devicePath)
	if err != nil {
		log.Fatalf("Failed to open MIDI device: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 3)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			log.Fatalf("Failed to read MIDI data: %v", err)
		}
		fmt.Print("MIDI Bytes: ")
		for i := range n {
			fmt.Printf("%02X ", buffer[i])
		}
		fmt.Println()

		if n == 3 && buffer[0] == 0x90 && buffer[1] == 0x32 && buffer[2] != 0x00 {
			fmt.Println("BUTTON = toggle")
			go func() {
				setLightStatus(myAuth, "toggle")
				time.Sleep(1 * time.Second)
			}()
		} else if n == 3 && buffer[0] == 0x90 && buffer[1] == 0x33 && buffer[2] != 0x00 {
			fmt.Println("BUTTON = on")
			go func() {
				setLightStatus(myAuth, "turn_on")
				time.Sleep(1 * time.Second)
			}()
		} else if n == 3 && buffer[0] == 0x90 && buffer[1] == 0x31 && buffer[2] != 0x00 {
			fmt.Println("BUTTON = off")
			go func() {
				setLightStatus(myAuth, "turn_off")
				time.Sleep(1 * time.Second)
			}()
		} else if n == 3 && buffer[0] == 0x90 && (buffer[1] >= 0x3C && buffer[1] <= 0x48) && buffer[2] != 0x00 {
			brightness_pct := int(float32(int(buffer[1])-60) * (100.0 / 12.0))
			fmt.Printf("BUTTON = brightness_pct: %d%%\n", brightness_pct)
			go func() {
				ChangeLightBrightnessByPct(myAuth, brightness_pct)
				time.Sleep(1 * time.Second)
			}()
		} else if n == 3 && buffer[0] == 0x90 && buffer[1] == 0x34 && buffer[2] != 0x00 {
			fmt.Printf("BUTTON = brightness increased by: %d\n", BRIGHTNESS_INCREMENT)
			go func() {
				IncrementLightBrigthness(myAuth, BRIGHTNESS_INCREMENT)
				time.Sleep(1 * time.Second)
			}()
		} else if n == 3 && buffer[0] == 0x90 && buffer[1] == 0x30 && buffer[2] != 0x00 {
			fmt.Printf("BUTTON = brightness decreased by: %d\n", BRIGHTNESS_INCREMENT)
			go func() {
				DecrementLightBrigthness(myAuth, BRIGHTNESS_INCREMENT)
				time.Sleep(1 * time.Second)
			}()
		}

	}

}

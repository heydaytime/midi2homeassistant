package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
)

type Auth struct {
	IP        string
	PORT      string
	ENDPOINT  string
	TOKEN     string
	ENTITY_ID string
}

type LightState struct {
	State      string `json:"state"`
	Attributes struct {
		Brightness uint8 `json:"brightness"`
	} `json:"attributes"`
}

func setLightStatus(auth Auth, status string) {

	URL := authToURL(auth, "services", "light", status)

	payload := map[string]string{"entity_id": auth.ENTITY_ID}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := createPostReq(URL, auth.TOKEN, jsonPayload)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

}

func ChangeLightBrightnessByPct(auth Auth, pct int) {
	URL := authToURL(auth, "services", "light", "turn_on")

	payload := map[string]string{"entity_id": auth.ENTITY_ID, "brightness_pct": fmt.Sprintf("%d", pct)}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := createPostReq(URL, auth.TOKEN, jsonPayload)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

func ChangeLightBrightnessByVal(auth Auth, val uint8) {

	URL := authToURL(auth, "services", "light", "turn_on")

	payload := map[string]string{"entity_id": auth.ENTITY_ID, "brightness": fmt.Sprintf("%d", val)}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := createPostReq(URL, auth.TOKEN, jsonPayload)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

func IncrementLightBrigthness(auth Auth, val int) {
	state, err := GetLightState(auth)
	if err != nil {
		log.Fatal(err)
	}

	newBrightness := uint8(math.Min(float64(state.Attributes.Brightness)+float64(val), 255))

	ChangeLightBrightnessByVal(auth, newBrightness)
}

func DecrementLightBrigthness(auth Auth, val int) {
	state, err := GetLightState(auth)
	if err != nil {
		log.Fatal(err)
	}

	newBrightness := uint8(math.Max(float64(state.Attributes.Brightness)-float64(val), 0))
	ChangeLightBrightnessByVal(auth, newBrightness)
}

func GetLightState(auth Auth) (LightState, error) {
	URL := authToURL(auth, "states", auth.ENTITY_ID)

	req, err := createGetReq(URL, auth.TOKEN)
	if err != nil {
		return LightState{}, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return LightState{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LightState{}, err
	}

	var lightState LightState

	err = json.Unmarshal(body, &lightState)
	if err != nil {
		return LightState{}, err
	}

	return lightState, nil
}

func createGetReq(url string, authToken string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func createPostReq(url string, authToken string, jsonPayload []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func authToURL(auth Auth, fn ...string) string {
	path := ""
	if len(fn) > 0 {
		path = "/" + strings.Join(fn, "/")
	}
	return fmt.Sprintf("%s:%s%s%s", auth.IP, auth.PORT, auth.ENDPOINT, path)
}

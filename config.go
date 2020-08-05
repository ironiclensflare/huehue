package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type HueConfig struct {
	IP          string       `json:"hue_ip"`
	User        string       `json:"hue_user"`
	DelayTimers []DelayTimer `json:"delay_timers"`
}

type DelayTimer struct {
	GroupID   int    `json:"group_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Delay     int
}

func GetConfig() (error, *HueConfig) {
	bytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/.huehue")
	if err != nil {
		log.Print("Failed to read file: ", err)
		return err, nil
	}

	var config HueConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Print("Failed to parse JSON: ", err)
		return err, nil
	}

	return nil, &config
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var lightTimerActivated bool
var turnLightOffTime time.Time
var config *HueConfig

func main() {
	_, config = GetConfig()
	spew.Dump(config)
	checkIfTheLightIsOn()
	timer := time.NewTicker(20 * time.Second)
	for range timer.C {
		checkIfTheLightIsOn()
	}
}

func rulesAreInEffect() bool {
	now := time.Now()
	if now.Hour() >= 12 || now.Hour() <= 8 {
		return true
	}
	return false
}

func buildRequestURL() string {
	return "http://" + config.IP + "/api/" + config.User
}

func checkIfTheLightIsOn() {
	if !rulesAreInEffect() {
		log.Print("Rules are not in effect. Skipping...")
		return
	}
	lightIsOn := isTheLightOn()
	if lightIsOn && lightTimerActivated {
		if turnLightOffTime.Sub(time.Now()) > 0 {
			log.Print("Light is on, timer is activated, not due to be turned off yet. Doing nothing...")
		} else {
			log.Print("Light is on, timer is activated and must be turned off!")
			turnTheLightOff()
		}
	} else if lightIsOn && !lightTimerActivated {
		log.Print("Light is on, timer is deactivated. Activating timer...")
		turnLightOffTime = time.Now().Add(5 * time.Minute)
		lightTimerActivated = true
	} else {
		log.Print("Light is off or timer is deactivated. Doing nothing...")
		lightTimerActivated = false
	}
}

func isTheLightOn() bool {
	url := buildRequestURL() + "/groups"
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	var objmap map[string]json.RawMessage
	json.Unmarshal(body, &objmap)
	var group Group
	json.Unmarshal(objmap["3"], &group)
	return group.State.Any_On
}

func turnTheLightOff() {
	url := buildRequestURL() + "/groups/3/action"
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPut, url, strings.NewReader(`{"on": false}`))
	client.Do(req)
}

type Group struct {
	Name   string
	Lights []string
	Type   string
	State  State
}

type State struct {
	All_On bool
	Any_On bool
}

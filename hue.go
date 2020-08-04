package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var lightTimerActivated bool
var turnLightOffTime time.Time

func main() {
	checkIfTheLightIsOn()
	timer := time.NewTicker(20 * time.Second)
	for range timer.C {
		checkIfTheLightIsOn()
	}
}

func rulesAreInEffect() bool {
	now := time.Now()
	if now.Hour() >= 22 || now.Hour() <= 8 {
		return true
	}
	return false
}

func buildRequestURL() string {
	ip := os.Getenv("HUE_IP")
	user := os.Getenv("HUE_USER")
	return "http://" + ip + "/api/" + user
}

func checkIfTheLightIsOn() {
	if !rulesAreInEffect() {
		fmt.Println("Rules are not in effect. Skipping...")
		return
	}
	lightIsOn := isTheLightOn()
	if lightIsOn && lightTimerActivated {
		if turnLightOffTime.Sub(time.Now()) > 0 {
			fmt.Println("Light is on, timer is activated, not due to be turned off yet. Doing nothing...")
		} else {
			fmt.Println("Light is on, timer is activated and must be turned off!")
			turnTheLightOff()
		}
	} else if lightIsOn && !lightTimerActivated {
		fmt.Println("Light is on, timer is deactivated. Activating timer...")
		turnLightOffTime = time.Now().Add(3 * time.Minute)
		lightTimerActivated = true
	} else {
		fmt.Println("Light is off or timer is deactivated. Doing nothing...")
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

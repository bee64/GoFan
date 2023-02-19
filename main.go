package main

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

const UP_INTERVAL = 5    // seconds between temp checks when getting hot
const DOWN_INTERVAL = 30 // seconds (in addition to UP) when cooling
const CYCLE = 10         // bigger cycle length allows more fine grained control

const MAX_TEMP = 80
const MIN_TEMP = 60 // temp when fan starts spinning
const TEMP_RANGE = MAX_TEMP / MIN_TEMP

var BACKGROUND = os.Getenv("BACKGROUND")

func getTemp() int {
	tempBytes, err := exec.Command("vcgencmd", "measure_temp").Output()
	must("read temp", err)
	// "temp=12.3'C" => 12
	temp := strings.Trim(string(tempBytes), "\n")
	temp = strings.ReplaceAll(temp, "temp=", "")
	temp = strings.ReplaceAll(temp, "'C", "")
	temp = strings.Split(temp, ".")[0]
	tempNum, err := strconv.Atoi(temp)
	must("convert temp to string", err)
	return tempNum
}

func getFanSpeed(temp int) uint32 {
	if temp <= MIN_TEMP {
		return 0
	} else if temp >= MAX_TEMP {
		return CYCLE
	} else {
		fanSpeed := (temp - MIN_TEMP) / TEMP_RANGE
		if fanSpeed == 1 {
			// min fanspeed of 20%, 10% is kinda ineffective
			return 2
		}
		return uint32(fanSpeed)
	}
}

func main() {
	checkRootUser()
	must("Open RPIO connection", rpio.Open())
	defer rpio.Close()
	pin := rpio.Pin(18)
	pin.Pwm()
	pin.Freq(25000 * CYCLE) // 25kHz for Noctua fan control
	pin.DutyCycle(0, CYCLE) // start with fan off

	lastFanSpeed := uint32(0)
	for {
		highTemp := getTemp()
		fanSpeed := getFanSpeed(highTemp)
		log("fan speed: " + strconv.FormatUint(uint64(fanSpeed*10), 10) + "%")
		log("temp:      " + strconv.Itoa(highTemp) + "C\n")

		// longer "cooldown hysteresis"
		if fanSpeed < lastFanSpeed {
			log("downward hysteresis wait")
			wait(DOWN_INTERVAL)
		}
		pin.DutyCycle(fanSpeed, CYCLE)
		lastFanSpeed = fanSpeed
		wait(UP_INTERVAL)
	}
}

// helper functions

func checkRootUser() {
	currentUser, err := user.Current()
	must("check current user", err)
	if currentUser.Username != "root" {
		println("PWM fan control needs to be run as root, try `sudo`")
		os.Exit(0)
	}
}

func wait(seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}

func log(message string) {
	if BACKGROUND == "" {
		println(message)
	}
}

// TODO alerting
func must(task string, err error) {
	if err != nil {
		panic(task + " error: " + err.Error())
	}
}

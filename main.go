package main

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

const INTERVAL = 5 // wait time between temp checks in seconds

func readCmd(cmd *exec.Cmd, name string) string {
	cmdOutput, err := cmd.Output()
	must(name, err)
	return strings.Trim(string(cmdOutput), "\n")
}

func getHighestTemp() int {
	// read & parse GPU temp: "temp=12.3'C" => 12
	gpuTemp := readCmd(exec.Command("vcgencmd", "measure_temp"), "read gpu temp")
	gpuTemp = strings.ReplaceAll(gpuTemp, "temp=", "")
	gpuTemp = strings.ReplaceAll(gpuTemp, "'C", "")
	gpuTemp = strings.Split(gpuTemp, ".")[0]
	gpuTempNum, err := strconv.Atoi(gpuTemp)
	must("convert gpu temp to string", err)

	// read & parse CPU temp: 12345 => 12
	// 	=> 12345 is 12.345'C
	cpuTemp := readCmd(exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp"), "read cpu temp")
	cpuTempNum, err := strconv.Atoi(cpuTemp)
	must("convert cpu temp to string", err)
	cpuTempNum /= 1000

	if gpuTempNum > cpuTempNum {
		return gpuTempNum
	} else {
		return cpuTempNum
	}
}

func main() {
	must("Open RPIO connection", rpio.Open())
	defer rpio.Close()
	pinPWM := rpio.Pin(18)

	highTemp := getHighestTemp()
	println("highest temp: " + strconv.Itoa(highTemp) + "C")

	// pinPWM.Output()
	// pinPWM.High()
	pinPWM.Pwm()
	pinPWM.Freq(25000 * 4) // 25kHz for Noctua fan control

	println("1")
	pinPWM.DutyCycle(1, 4)
	wait()
	println("2")
	pinPWM.DutyCycle(2, 4)
	wait()
	println("4")
	pinPWM.DutyCycle(4, 4)
	wait()
	println("off")
	pinPWM.DutyCycle(0, 4)
	wait()
	// println("0")
}

func wait() {
	time.Sleep(time.Second * INTERVAL)
}

// TODO alerting
func must(task string, err error) {
	if err != nil {
		panic(task + " error: " + err.Error())
	}
}

// pins used by the fan that don't need programmatic control
// pin5V  := rpio.Pin(4) // could use pin 2 as well
// pin3V  := rpio.Pin(17) // could use pin 1 as well
// pinGND := rpio.Pin(6)
// pinPWM := rpio.Pin(12) // GPIO 18

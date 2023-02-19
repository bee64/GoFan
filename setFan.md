# adds "sudo SPEED=2 ./GoFan"

add this to `main` after pin setup, and run with `sudo SPEED=2 ./GoFan`

```go
if os.Getenv("SPEED") != "" {
	speed, err := strconv.ParseUint(os.Getenv("SPEED"), 10, 64)
	must("one-time speed", err)
	log("set one time speed: " + os.Getenv("SPEED"))
	pinPWM.DutyCycle(uint32(speed), CYCLE)
	os.Exit(0)
}
```

this sets the PWM speed and exits, which was useful for me in figuring out a good floor for the fan speed. could also use it to set "always run at 50%", but that's kinda not the point of this so i left it out of the code
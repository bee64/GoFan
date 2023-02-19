# GoFan
a PWM fan controller built for the Raspberry Pi 4 & Noctua NF-A4x10 5V PWM

## building
i'm running ubuntu so I build with `GOOS=linux GOARCH=arm go build .`.

## running

`sudo ./GoFan` 
> [root is required](https://github.com/stianeikeland/go-rpio#using-without-root) to use PWM pins

### background mode

run with `sudo BACKGROUND="anything but empty string" ./GoFan &` to run as a background process

## ENV
`BACKGROUND` changes logging behavior. If this var is present, logs are not printed to the console and are instead ~written to logfile TODOLOL~ yeeted to the void

### References

DriftKingTw's [pi PWM fan](https://blog.driftking.tw/en/2019/11/Using-Raspberry-Pi-to-Control-a-PWM-Fan-and-Monitor-its-Speed/) was the main source I used in building this, with the main difference being the use of Golang instead of Python to write the fan controller

Check out their guide for pin setups & general PWM fan theory!
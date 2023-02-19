# GoFan
a PWM fan controller built for the Raspberry Pi 4 & Noctua NF-A4x10 5V PWM

## building
i'm running ubuntu so I build with `GOOS=linux GOARCH=arm go build .`

## notes
the gpio library requires root to to use PWM mode, so run with `sudo`
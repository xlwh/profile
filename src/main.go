package main

import (
	"profile"
	"time"
)

func main() {
	p := profile.GetProfile("./profile.dump", 10)
	p.Start()
	time.Sleep(time.Hour * 1)
	p.Stop()
}

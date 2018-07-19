package main

import (
	"fmt"
	"profile"
)

func main() {
	fmt.Println("Start profile")
	p := profile.GetProfile("./profile.dump", 10)

	p.Start()
//	time.Sleep(time.Millisecond * 300000)
	p.Stop()

	fmt.Println("done")
}

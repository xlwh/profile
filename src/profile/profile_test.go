package profile

import (
	"testing"
	"time"
	"fmt"
)

func TestProfile(t *testing.T) {
	fmt.Println("Start profile")
	p := GetProfile("./profile.dump", 10)
	p.Start()
	time.Sleep(time.Millisecond * 300000)
	p.Stop()

	fmt.Println("done")
}

package main

import (
	"fmt"

	"github.com/go-olive/olive/foundation/biliup"
)

func main() {
	err := biliup.New(biliup.Config{
		CookieFilepath:    "/Users/xxx/cookies.json",
		VideoFilepath:     `/video.mp4`,
		Threads:           2,
		MaxBytesPerSecond: 0,
	}).Upload()
	if err != nil {
		fmt.Println(err)
	}
}

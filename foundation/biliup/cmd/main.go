package main

import (
	"fmt"

	"github.com/go-olive/olive/foundation/biliup"
)

func main() {
	err := biliup.New(biliup.Config{
		CookieFilepath:    "/cookies.json",
		VideoFilepath:     `/test.flv`,
		Threads:           2,
		MaxBytesPerSecond: 2097152,
	}).Upload()
	if err != nil {
		fmt.Println(err)
	}
}

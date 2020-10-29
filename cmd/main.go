package main

import (
	"fmt"
	"os"

	"github.com/classAndrew/dynmapgo/pkg/client"
)

func main() {
	url := os.Args[1]
	fmt.Println("Connecting to " + url)
	if err := client.Connect(url); err != nil {
		fmt.Println("Unable to connect\n\n" + err.Error())
	}
}

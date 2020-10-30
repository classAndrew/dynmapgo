package main

import (
	"fmt"
	"os"

	"github.com/classAndrew/dynmapgo/pkg/client"
)

func main() {
	url := os.Args[1]
	fmt.Println("Connecting to " + url)
	cl := client.Client{URL: url}
	if err := cl.Connect(); err != nil {
		fmt.Println("Unable to connect\n\n" + err.Error())
	}
	cl.DownloadMap(25, 25, 1, 0, 0)
	cl.CompositeLeaflets(25, 25, 1, 0, 0)
}

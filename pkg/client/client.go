package client

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const dynmapBoilerplate string = "<meta name=\"description\" content=\"Minecraft Dynamic Map\" />"

// in the form tiles/world/flat/-1_1/zzzzz_-32_32.jpg
const imgworld string = "tiles/world/flat/"

// Client to the server
type Client struct {
	// URL to connect to
	URL string
}

// Connect no real connection, just check if server is responsive
func (cl *Client) Connect() error {
	res, err := http.Get(cl.URL)
	if res == nil || int(res.ContentLength) < len(dynmapBoilerplate) {
		return errors.New("Not a dynmap server or server not found- Please double check that you've entered things correctly")
	}
	content := make([]byte, res.ContentLength)
	res.Body.Read(content)
	if err != nil {
		return err
	} else if !strings.Contains(string(content), dynmapBoilerplate) {
		return errors.New("Not a dynmap server- Please double check you that you entered the host correctly")
	}
	return nil
}

// DownloadMap will download all the leaflets in a rectangular region with height x width leaflets. The scale is a power of 4
func (cl *Client) DownloadMap(height int, width int, scale int, x int, y int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			coordname := strconv.Itoa((j-width/2+x)*32) + "_" + strconv.Itoa((i-height/2+y)*32) + ".jpg"
			res, _ := http.Get(cl.URL + imgworld + "0_0" + "/zzzzz_" + coordname)
			img, _ := ioutil.ReadAll(res.Body)
			ioutil.WriteFile("./leaflets/"+coordname, img, 0644)
			fmt.Printf("Done with %d images\n", i*height+j)
		}
	}
}

// DownloadMapAsync will download all the leaflets asynchronously
func (cl *Client) DownloadMapAsync(height int, width int, scale int, x int, y int) {
	var wg sync.WaitGroup
	count := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			coordname := strconv.Itoa((j-width/2+x)*32) + "_" + strconv.Itoa((i-height/2+y)*32) + ".jpg"
			// 10 so ratelimit doesn't freakout
			if count == 40 {
				wg.Wait()
				count = 0
			}
			wg.Add(1)
			count++
			go func(coordname string, i, j int) {
				res, _ := http.Get(cl.URL + imgworld + "0_0" + "/zzzzz_" + coordname)
				img, _ := ioutil.ReadAll(res.Body)
				ioutil.WriteFile("./leaflets/"+coordname, img, 0644)
				fmt.Printf("Done with %d images\n", i*height+j)
				wg.Done()
			}(coordname, i, j)
		}
	}
	wg.Wait()
}

// CompositeLeaflets composites all the leaflets into one big image
func (cl *Client) CompositeLeaflets(height int, width int, scale int, x int, y int) {
	im := image.NewRGBA(image.Rect(0, 0, width*128, height*128))
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			file, err1 := os.Open("./leaflets/" + strconv.Itoa((j-width/2+x)*32) + "_" + strconv.Itoa((i-height/2+y)*32) + ".jpg")
			if err1 != nil {
				fmt.Println(err1.Error())
			}
			defer file.Close()
			src, err := jpeg.Decode(file)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			draw.Draw(im, image.Rect(j*128, (height-i)*128, (j+1)*128, (height-i-1)*128), src, image.Pt(0, 0), draw.Over)
		}
	}
	// ioutil.WriteFile("big.png", []byte(""), 0644)
	im.Set(1, 1, color.RGBA{20, 20, 20, 100})
	big, err := os.OpenFile("big.jpg", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer big.Close()
	err = jpeg.Encode(big, im, &jpeg.Options{Quality: 100})
	if err != nil {
		fmt.Println(err)
	}
}

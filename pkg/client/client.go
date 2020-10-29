package client

import (
	"errors"
	"net/http"
	"strings"
)

const dynmapBoilerplate string = "<meta name=\"description\" content=\"Minecraft Dynamic Map\" />"

// Connect no real connection, just check if server is responsive
func Connect(url string) error {
	res, err := http.Get(url)
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

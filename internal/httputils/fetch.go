package httputils

import (
	"fmt"
	"io"
	"net/http"
)

// Fetch will download a url to a Writer.
func Fetch(url string, wri io.Writer) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching '%s' (status code: %d)", url, resp.StatusCode)
	}

	// Write the body to file
	_, err = io.Copy(wri, resp.Body)
	return err
}

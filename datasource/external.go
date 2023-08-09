package datasource

import (
	"log"
	"net/http"
)

// Get makes a HTTP request to the url provided
func Get(url string, resp interface{}) error {
	log.Printf("INFO: making a request with URL: %s\n", url)

	// make the request
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make http request to url: %s. Error: %v\n", url, err)
	}

	// read all the response body
	
	return nil
}

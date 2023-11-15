package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func getPageHtml(url string, errChan chan<- error) (node *html.Node) {
	res, err := http.Get(url)
	if err != nil {
		errChan <- fmt.Errorf("error fetching page %s: %v", url, err)
		return nil
	}
	defer res.Body.Close()

	node, err = html.Parse(res.Body)
	if err != nil {
		errChan <- fmt.Errorf("error parsing page %s: %v", url, err)
		return nil
	}

	return
}

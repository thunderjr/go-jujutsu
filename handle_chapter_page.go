package main

import (
	"fmt"
	"strconv"
	"sync"
)

func handleChapterPage(prefix, url string, errChan chan<- error, fileChan chan *File, wg *sync.WaitGroup) {
	defer wg.Done()

	mainNode := getPageHtml(url, errChan)
	if mainNode == nil {
		errChan <- fmt.Errorf("handleChapterPage (%s): mainNode is nil", url)
		return
	}

	imgs := queryAll(mainNode, "main > article > div > div > p img")
	// imgs := queryAll(mainNode, "main > article > div > div > div > a img")

	for pageNum, page := range imgs {
		for _, attr := range page.Attr {
			if attr.Key == "src" {
				key := strconv.Itoa(pageNum + 1)
				fileChan <- &File{
					Url:  attr.Val,
					Name: fmt.Sprintf("%s/%s", prefix, key),
				}
			}
		}
	}
}

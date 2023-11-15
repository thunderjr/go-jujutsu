package main

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

const URL = "https://ww7.jujmanga.com/"

func main() {
	wg := &sync.WaitGroup{}

	fileChan := make(chan *File)
	errChan := make(chan error)

	defer close(fileChan)
	defer close(errChan)

	go handleFiles(fileChan, errChan, wg)
	go handleErrors(errChan)

	homeNode := getPageHtml(URL, errChan)
	if homeNode == nil {
		return
	}

	chapterLinks := queryAll(homeNode, "#ceo_latest_comics_widget-3 > ul > li > a")
	if len(chapterLinks) == 0 {
		errChan <- fmt.Errorf("no chapter links found")
		return
	}

	for i, link := range chapterLinks {
		for _, attr := range link.Attr {
			if attr.Key == "href" {
				linkChunks := strings.Split(link.FirstChild.Data, " ")
				chapter := linkChunks[len(linkChunks)-1]

				fileName := fmt.Sprintf("./output/%s", chapter)

				wg.Add(1)
				go handleChapterPage(fileName, attr.Val, errChan, fileChan, wg)

				if math.Mod(float64(i), 10) == 0 {
					time.Sleep(2 * time.Second)
				}
			}
		}
	}

	wg.Wait()
}

func handleErrors(errChan <-chan error) {
	for err := range errChan {
		fmt.Printf("[ERROR] %+v\n", err)
	}
}

package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	Url  string
	Name string
}

func handleFiles(fileChan <-chan *File, errChan chan<- error, wg *sync.WaitGroup) {
	for file := range fileChan {
		wg.Add(1)
		go downloadFile(file, errChan, wg)
	}
}

func downloadFile(file *File, errChan chan<- error, wg *sync.WaitGroup) {
	err := download(file.Url, file.Name)
	if err != nil {
		errChan <- fmt.Errorf("downloadFile (%s): %w", file.Url, err)
	}

	wg.Done()
}

func download(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	ext, err := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("failed to get file extension: %w", err)
	}

	filename += ext[len(ext)-1]
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

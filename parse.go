package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ParseError struct {
	Description string
	error       error
}

func (parseError *ParseError) Error() string {
	return fmt.Sprintf("Something went wrong: %s", parseError.Description)
}

func handleRequestError(res *http.Response, err error) error {
	if err != nil {
		return &ParseError{error: err, Description: "Not able to connect the server"}
	}
	if res.StatusCode != 200 {
		return &ParseError{error: nil, Description: "Thread not found"}
	}
	return nil
}

func logAndSkipError(err error) {
	f, _ := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.Printf("Skipped image: %s", err.Error())
}

func parse(url string) (*http.Response, error) {
	res, err := http.Get(url)

	return res, handleRequestError(res, err)
}

func downloadFile(filePath string, url string) {
	response, err := http.Get(url)
	if err != nil {
		logAndSkipError(err)
		return

	}

	if _, err := os.Stat(filePath); err == nil {
		return
	}

	data, _ := ioutil.ReadAll(response.Body)

	response.Body.Close()
	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		logAndSkipError(err)
	}

}

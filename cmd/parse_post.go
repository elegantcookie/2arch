package cmd

import (
	"fmt"
	"io"
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
	log.Printf("Skipped file: %s reload started...", err.Error())
}

func parse(url string) (*http.Response, error) {
	res, err := http.Get(url)

	return res, handleRequestError(res, err)
}

func createFile(filePath string) (*os.File, error) {
	output, err := os.Create(filePath)
	if err != nil {
		logAndSkipError(err)
		return output, err
	}
	return output, nil
}

func downloadFile(filePath string, url string) bool {
	//if _, err := os.Stat(filePath); err == nil {
	//	return false
	//}

	output, err := createFile(filePath)
	if err != nil {
		logAndSkipError(err)
		return false
	}
	defer output.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, http.NoBody)
	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		logAndSkipError(err)
		return false
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)

	if err != nil {
		fmt.Println(ioutil.ReadAll(response.Body))
		logAndSkipError(err)
		return false
	}
	return true

}

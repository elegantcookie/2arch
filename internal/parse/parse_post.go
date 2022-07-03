package parse

import (
	"2arch/pkg/logging"
	"fmt"
	"io"
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

func parse(url string) (*http.Response, error) {
	res, err := http.Get(url)
	return res, handleRequestError(res, err)
}

// Создаёт файл по пути filePath
func createFile(filePath string) (*os.File, error) {
	output, err := os.Create(filePath)
	if err != nil {
		logging.LogAndSkipError(err)
		return output, err
	}
	return output, nil
}

func downloadFile(filePath string, url string) bool {
	output, err := createFile(filePath)
	if err != nil {
		logging.LogAndSkipError(err)
		return false
	}
	defer output.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, http.NoBody)
	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		logging.LogAndSkipError(err)
		return false
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)

	if err != nil {
		logging.LogAndSkipError(err)
		return false
	}
	return true

}

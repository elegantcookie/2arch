package main

import (
	"fmt"
	"net/http"
)

type ParseError struct {
	Description string
	error       error
}

func (parseError *ParseError) Error() string {
	return fmt.Sprintf("Something went wrong: %s", parseError.Description)
}

func parse(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, &ParseError{error: err, Description: "Not able to connect the server"}
	}
	if res.StatusCode != 200 {
		return nil, &ParseError{error: nil, Description: "Thread not found"}
	}

	return res, nil
}

package main

import (
	"fmt"
	"net/http"
)

type ParseError struct {
	error error
}

func (parseError *ParseError) Error() string {
	return fmt.Sprintf("Something went wrong: %s", parseError.error.Error())
}

func parse(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, &ParseError{error: err}
	}

	return res, nil
}

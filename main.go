package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Post struct {
	Comment string
	Files   interface{}
}

type PostsArray struct {
	Posts []Post
}

type ThreadInfo struct {
	PostCount map[string]int
	Threads   []PostsArray
}

type UrlMatchError struct {
	url   string
	error error
}

func (matchError *UrlMatchError) Error() string {
	if matchError.error != nil {
		return matchError.error.Error()
	}
	return fmt.Sprintf("Wrong url: %s\n", matchError.url)
}

func matchCase(url string) (bool, error) {
	exrp := "https?:\\/\\/2ch\\.(hk|life)\\/\\w+\\/res\\/\\d+\\.html"
	matched, err := regexp.Match(exrp, []byte(url))
	if err != nil {
		return false, &UrlMatchError{url: url, error: err}
	}
	if matched {
		return true, nil
	} else {
		return false, &UrlMatchError{url: url, error: nil}
	}

}

func handleError(error error) error {
	if error != nil {
		fmt.Printf("An error has occured: %s", error.Error())
		os.Exit(-1)
	}
	return nil
}

func main() {
	htmlUrl := "https://2ch.hk/b/res/269256093.html"

	matched, err := matchCase(htmlUrl)
	handleError(err)

	if matched {
		fmt.Println("Good url")
	} else {
		fmt.Println("Wrong url")
		return
	}

	m1 := regexp.MustCompile(`html`)

	jsonUrl := m1.ReplaceAllString(htmlUrl, "json")

	res, err := parse(jsonUrl)
	handleError(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	handleError(err)

	var threadInfo ThreadInfo

	err = json.Unmarshal(body, &threadInfo)
	handleError(err)

	fmt.Println(threadInfo.Threads[0].Posts[2].Files)
}

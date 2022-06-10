package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type Post struct {
	Num       int
	Name      string
	Timestamp int
	Comment   string
	Files     interface{}
}

type PostsArray struct {
	Posts []Post
}

type ThreadInfo struct {
	Posts_Count int
	Board       string
	Threads     []PostsArray
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

func renderPostHtml(postClassName string, args ...interface{}) string {
	postHtml := fmt.Sprintf("<div id=\"thread-%d\" class=\""+
		postClassName+
		"\">"+
		"<p>%s %s <a href=\"%s/res/#%d.html\">№%d</a>"+
		"</p>"+
		"<p>%s</p>"+
		"</div><br>", args...)
	return postHtml
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

	postNum := threadInfo.Threads[0].Posts[0].Num
	postText := threadInfo.Threads[0].Posts[0].Comment

	// THERE MUST BE A TEMPLATE ENGINE
	htmlFile := fmt.Sprintf("<!DOCTYPE html><html><head>2ch-archiver</head><body><h1>Start of file</h1><div id=\"thread-%d\" class=\"thread\">", postNum)

	for i := 0; i < threadInfo.Posts_Count; i++ {
		postNum = threadInfo.Threads[0].Posts[i].Num
		postText = threadInfo.Threads[0].Posts[i].Comment
		datetime := timestampToDatetime(threadInfo.Threads[0].Posts[i].Timestamp)

		if i == 0 {
			htmlFile += renderPostHtml("thread__oppost", postNum, datetime.Date, datetime.Time, threadInfo.Board, postNum, postNum, postText)

		} else {
			htmlFile += renderPostHtml("thread__post", postNum, datetime.Date, datetime.Time, threadInfo.Board, postNum, postNum, postText)
		}
	}
	htmlFile += "</dib></body></html>"
	ioutil.WriteFile(fmt.Sprintf("%d.html", threadInfo.Threads[0].Posts[0].Num), []byte(htmlFile), 0600)
}

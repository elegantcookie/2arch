package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type PostFile struct {
	DisplayName string
	FullName    string
	Path        string
	Size        int
	Thumbnail   string
	Tn_Height   int
	Tn_Width    int
}

//
//type PostFiles struct {
//	Files []PostFile
//}

type Post struct {
	Num       int
	Name      string
	Timestamp int
	Comment   string
	Files     []PostFile
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

func isMatch(str string, expr string) bool {
	matched, _ := regexp.Match(expr, []byte(str))
	if matched {
		return true
	} else {
		return false
	}
}

func handleError(error error) error {
	if error != nil {
		panic(error.Error())

	}
	return nil
}

func renderPostHtml(postClassName string, args ...interface{}) string {
	postHtml := fmt.Sprintf("<div id=\"thread-%d\" class=\""+
		postClassName+
		"\">"+
		"<p>%s %s <a href=\"#%d\">№%d</a>"+
		"</p>"+
		"<p>%s</p>"+
		"</div><br>", args...)
	return postHtml
}

func main() {
	htmlUrl := "https://2ch.hk/b/res/269314054.html"

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

	rootPath, _ := os.Getwd()

	threadNum := threadInfo.Threads[0].Posts[0].Num
	path := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d/", threadNum))

	os.MkdirAll(path, os.ModePerm)

	// THERE MUST BE A TEMPLATE ENGINE
	htmlFile := fmt.Sprintf("<!DOCTYPE html><html><head>2ch-archiver</head><body><h1>Start of file</h1><div id=\"thread-%d\" class=\"thread\">", postNum)

	m2, _ := regexp.Compile(`/\w/res.+#`)

	var wg sync.WaitGroup

	for i := 0; i < threadInfo.Posts_Count; i++ {
		postNum = threadInfo.Threads[0].Posts[i].Num
		postText = threadInfo.Threads[0].Posts[i].Comment
		postText = m2.ReplaceAllString(postText, fmt.Sprintf("%s/%d.html#", path, threadNum))
		datetime := timestampToDatetime(threadInfo.Threads[0].Posts[i].Timestamp)

		files := threadInfo.Threads[0].Posts[i].Files

		filesNum := uint64(len(files))
		if filesNum != 0 {
			for j := range files {
				fmt.Println(j)
				fmt.Println(files[j].Path)
				fmt.Println(fmt.Sprintf("2ch.hk/%s", files[j].Path))

				if isMatch(files[j].Path, "sticker") {
					continue
				}

				wg.Add(1)
				j := j
				go func() {
					ok := false
					for !ok {
						ok = downloadFile(fmt.Sprintf("%s/%s", path, files[j].FullName), fmt.Sprintf("http://2ch.hk%s", files[j].Path))
					}
					defer wg.Done()
				}()
			}
		}

		if i == 0 {
			htmlFile += renderPostHtml("thread__oppost", postNum, datetime.Date, datetime.Time, postNum, postNum, postText)

		} else {
			htmlFile += renderPostHtml("thread__post", postNum, datetime.Date, datetime.Time, postNum, postNum, postText)
		}
	}
	htmlFile += "</div></body></html>"

	fmt.Println(path)
	ioutil.WriteFile(fmt.Sprintf("threads/thread_%d/%d.html", threadNum, threadNum), []byte(htmlFile), 0600)
	wg.Wait()
}

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)

type ViewData struct {
	RootPath string
	Posts    []Post
}

type PostFile struct {
	Name               string
	DisplayName        string
	Duration           string
	FullName           string
	Path               string
	LocalPath          string
	LocalThumbnailPath string
	Md5                string
	Size               int
	Type               int
	Thumbnail          string
	Height             int
	Width              int
	Tn_Height          int
	Tn_Width           int
}

//
//type PostFiles struct {
//	Files []PostFile
//}

type Post struct {
	Name      string
	Subject   string
	Date      string
	Num       int
	Number    int
	Timestamp int
	Comment   string
	ImgAmount int
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

func downloadJson(htmlUrl string) {
	start := time.Now()
	matched, err := matchCase(htmlUrl)
	handleError(err)

	if matched {
		fmt.Println("Good url")
	} else {
		fmt.Println("Wrong url")
		return
	}

	m1 := regexp.MustCompile(`html`) // html url -> json url

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
	filesPath := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d/", threadNum))

	os.MkdirAll(filesPath, os.ModePerm)

	var wg sync.WaitGroup
	var mtx sync.Mutex
	postsMap := make(map[string]string)
	for i := 0; i < threadInfo.Posts_Count; i++ {
		postNum = threadInfo.Threads[0].Posts[i].Num
		postText = threadInfo.Threads[0].Posts[i].Comment

		wg.Add(1)
		go func() {
			mtx.Lock()
			postsMap[fmt.Sprintf("post-%d", postNum)] = postText
			defer mtx.Unlock()
			defer wg.Done()
		}()
	}
	output, err := createFile(filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d/%d.json", threadNum, threadNum)))
	handleError(err)
	defer output.Close()
	marshalledMap, err := json.Marshal(postsMap)
	handleError(err)

	strMap := string(marshalledMap)
	_, err = io.Copy(output, strings.NewReader(strMap))
	if err != nil {
		logAndSkipError(err)
		return
	}

	wg.Wait()
	defer log.Printf("elapsed time: %s", time.Since(start))
}

func downloadHtml(htmlUrl string) {

	start := time.Now()
	matched, err := matchCase(htmlUrl)
	handleError(err)

	if matched {
		fmt.Println("Good url")
	} else {
		fmt.Println("Wrong url")
		return
	}

	m1 := regexp.MustCompile(`html`) // html url -> json url
	m2 := regexp.MustCompile(`\..+`) // thumbnail url -> thumbnail.jpg
	m3 := regexp.MustCompile("/\\w+/\\w+/\\w+.\\w+#")

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

	filesPath := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d/files", threadNum))
	htmlFilePath := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d", threadNum))

	os.MkdirAll(filesPath, os.ModePerm)
	os.Mkdir(filesPath+"/thumbnails", os.ModePerm)

	var wg sync.WaitGroup

	var postsArray []Post
	for i := 0; i < threadInfo.Posts_Count; i++ {
		postNum = threadInfo.Threads[0].Posts[i].Num
		postText = threadInfo.Threads[0].Posts[i].Comment
		rPostNum := threadInfo.Threads[0].Posts[i].Number
		timestamp := threadInfo.Threads[0].Posts[i].Timestamp
		date := threadInfo.Threads[0].Posts[i].Date
		subject := threadInfo.Threads[0].Posts[i].Subject
		name := threadInfo.Threads[0].Posts[i].Name
		files := threadInfo.Threads[0].Posts[i].Files

		postText = m3.ReplaceAllString(postText, "#post-")

		for k := range files {
			files[k].LocalPath = fmt.Sprintf("files/%s", files[k].Name)
			thumbnailUrl := m2.ReplaceAllString(files[k].Name, ".jpg")
			files[k].LocalThumbnailPath = fmt.Sprintf("files/thumbnails/%s", thumbnailUrl)
			//fmt.Printf("thumbnail url: %s\n", files[k].Thumbnail)
			//fmt.Printf("local thumbnail url: %s\n", files[k].LocalThumbnailPath)
		}

		var imgAmount int
		if len(files) != 0 {
			imgAmount = len(files)
		} else {
			imgAmount = 0
		}

		postsArray = append(postsArray, Post{
			ImgAmount: imgAmount,
			Name:      name,
			Subject:   subject,
			Date:      date,
			Num:       postNum,
			Number:    rPostNum,
			Comment:   postText,
			Timestamp: timestamp,
			Files:     files,
		})

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
						if files[j].Type == 10 {
							for !ok {
								thumbnailUrl := m2.ReplaceAllString(files[j].Name, ".jpg")
								_localPath := fmt.Sprintf("%s/thumbnails/%s", filesPath, thumbnailUrl)
								_webPath := fmt.Sprintf("http://2ch.hk%s", files[j].Thumbnail)
								//fmt.Printf("localPath: %s\nweb_path: %s\n", _localPath, _webPath)
								ok = downloadFile(_localPath, _webPath)
							}
						}
						ok = false
						_localPath := fmt.Sprintf("%s/%s", filesPath, files[j].Name)
						_webPath := fmt.Sprintf("http://2ch.hk%s", files[j].Path)
						//fmt.Printf("localPath: %s\nweb_path: %s\n", _localPath, _webPath)
						ok = downloadFile(_localPath, _webPath)
					}
					defer wg.Done()
				}()
			}
		}

		data := ViewData{
			RootPath: rootPath,
			Posts:    postsArray,
		}
		tmpl, _ := template.ParseFiles("index.html")
		f, _ := os.Create("parsed.html")

		tmpl.Execute(f, data)
		f.Close()
		pathToFolder := fmt.Sprintf("%s/%d.html", htmlFilePath, threadNum)
		//fmt.Println(pathToFolder)
		f, _ = os.Create(pathToFolder)
		tmpl.Execute(f, data)
		f.Close()

	}

	fmt.Println(htmlFilePath)
	wg.Wait()
	defer log.Printf("elapsed time: %s", time.Since(start))
}

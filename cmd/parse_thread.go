package cmd

import (
	"2arch/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
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

// SimplifiedFile - упрощенная структура PostFile, чтобы хранить массив меньшего объема
type SimplifiedFile struct {
	Filename string
	Size     int
	Type     int
	Height   int
	Width    int
}

// ViewData - данные для заполнения шаблона
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
	IsVideo            bool
}

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
	if len(matchError.url) == 0 {
		return fmt.Sprint("не указан адрес треда")
	}
	return fmt.Sprint("некорректно указан адрес треда")
}

// Проверка ссылки на тред на корректность
func matchCase(url string) (bool, error) {
	exrp := "https?:\\/\\/2ch\\.hk\\/\\w+\\/res\\/\\d+\\.html"
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

// Проверка строки str на соответствие рег выражению expr
func isMatch(str string, expr string) bool {
	matched, _ := regexp.Match(expr, []byte(str))
	if matched {
		return true
	} else {
		return false
	}
}

// Простой обработчик ошибок
func handleError(error error) error {
	if error != nil {
		fmt.Printf("Произошла ошибка: %s", error.Error())
		os.Exit(-1)
	}
	return nil
}

type Thread struct {
	threadInfo ThreadInfo
	rootPath   string
}

// Проверяет ссылку на тред, если она действительна, то возвращает структуру треда
func setupThread(htmlUrl string) *Thread {
	matched, err := matchCase(htmlUrl)
	handleError(err)

	if matched {
		logging.LogMessage("Good url")
	} else {
		logging.LogMessage("Wrong url")
		handleError(fmt.Errorf("wrong url"))
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

	rootPath, _ := os.Getwd()

	return &Thread{
		threadInfo: threadInfo,
		rootPath:   rootPath,
	}
}

func isVideo(fileType int) bool {
	if fileType == 10 || fileType == 6 {
		return true
	}
	return false
}

func isImage(fileType int) bool {
	if fileType == 1 || fileType == 2 {
		return true
	}
	return false
}

// Скачивает страницу и создает json файл с текстом ответов
func downloadJson(htmlUrl string) {
	start := time.Now()

	logger := logging.Init("logs.txt")
	defer logger.Close()

	thread := setupThread(htmlUrl)
	threadInfo := thread.threadInfo
	rootPath := thread.rootPath

	postNum := threadInfo.Threads[0].Posts[0].Num
	postText := threadInfo.Threads[0].Posts[0].Comment

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
		logging.LogAndSkipError(err)
		return
	}

	wg.Wait()
	defer log.Println(fmt.Sprintf("Тред №%d скачан за: %s", threadNum, time.Since(start)))
}

// Скачивает страницу и заполняет её шаблон
func downloadHtml(htmlUrl string, flags Flags) {
	start := time.Now()

	logger := logging.Init("logs.txt")
	defer logger.Close()

	m2 := regexp.MustCompile(`\..+`) // thumbnail url -> thumbnail.jpg
	m3 := regexp.MustCompile("/\\w+/\\w+/\\w+.\\w+#")

	thread := setupThread(htmlUrl)
	threadInfo := thread.threadInfo
	rootPath := thread.rootPath

	var fileSet FileSet
	var filenameSet FilenameSet

	postNum := threadInfo.Threads[0].Posts[0].Num
	postText := threadInfo.Threads[0].Posts[0].Comment

	threadNum := threadInfo.Threads[0].Posts[0].Num

	filesPath := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d/files", threadNum))
	htmlFilePath := filepath.Join(rootPath, fmt.Sprintf("/threads/thread_%d", threadNum))

	os.MkdirAll(filesPath, os.ModePerm)
	os.Mkdir(filesPath+"/thumbnails", os.ModePerm)

	var wg sync.WaitGroup

	var postsArray []Post

	bar := pb.StartNew(0)
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
			thumbnailUrl := m2.ReplaceAllString(files[k].FullName, ".jpg")
			files[k].LocalThumbnailPath = fmt.Sprintf("files/thumbnails/%s", thumbnailUrl)
			files[k].IsVideo = isVideo(files[k].Type)
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
				// Если скачиваются только изображения, то остальные файлы пропускаются
				if flags.imagesOnly {
					if !isImage(files[j].Type) {
						continue
					}
				}

				// Если скачиваются только видео, то остальные файлы пропускаются
				if flags.videosOnly {
					if !isVideo(files[j].Type) {
						continue
					}
				}

				sFile := SimplifiedFile{
					Filename: files[j].FullName,
					Size:     files[j].Size,
					Type:     files[j].Type,
					Height:   files[j].Height,
					Width:    files[j].Width,
				}

				// Если файл во множестве fileSet, то у files[j] ссылка будет указывать на адрес этого файла
				if foundFile := fileSet.Find(sFile); foundFile != nil {
					filename := foundFile.Filename
					if !strings.Contains(filename, ".") {
						filename += ".jpg"
					}
					files[j].LocalPath = fmt.Sprintf("%s/%s", filesPath, filename)
					continue
				}

				// Если у файла одинаковое имя с другими, но он уникальный, то ему присваивается случайное имя
				if filenameSet.Contains(files[j].FullName) {
					files[j].FullName = GenerateFilename()
				}

				// Обновление имени в структуре sFile
				sFile.Filename = files[j].FullName

				// Добавление нового элемента во множества
				fileSet.Add(sFile)
				filenameSet.Add(files[j].FullName)

				bar.AddTotal(1)

				logging.LogMessage(fmt.Sprintf("download started for 2ch.hk/%s", files[j].Path))

				// Стикеры не скачиваются
				if isMatch(files[j].Path, "sticker") {
					continue
				}

				wg.Add(1)
				j := j
				// Нуждается в рефакторинге
				go func() {
					ok := false
					for !ok {
						if isVideo(files[j].Type) {
							// Скачивание превью для видео
							for !ok {
								thumbnailUrl := m2.ReplaceAllString(files[j].FullName, ".jpg")
								_localPath := fmt.Sprintf("%s/thumbnails/%s", filesPath, thumbnailUrl)

								_webPath := fmt.Sprintf("http://2ch.hk%s", files[j].Thumbnail)
								ok = downloadFile(_localPath, _webPath)
							}
						}
						ok = false
						_localPath := fmt.Sprintf("%s/%s", filesPath, files[j].FullName)
						files[j].LocalPath = _localPath

						_webPath := fmt.Sprintf("http://2ch.hk%s", files[j].Path)
						ok = downloadFile(_localPath, _webPath)
					}
					defer bar.Increment()
					defer wg.Done()
				}()
			}
		}

		data := ViewData{
			RootPath: rootPath,
			Posts:    postsArray,
		}
		tmpl, _ := template.ParseFiles("index.html")

		pathToFolder := fmt.Sprintf("%s/%d.html", htmlFilePath, threadNum)
		f, _ := os.Create(pathToFolder)
		tmpl.Execute(f, data)
		f.Close()

	}
	logMessage(fmt.Sprintf("files saved at %s", htmlFilePath))
	wg.Wait()
	bar.Finish()
	defer logMessage(fmt.Sprintf("Тред №%d скачан за: %s", threadNum, time.Since(start)))
}

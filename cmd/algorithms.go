package cmd

import (
	"math/rand"
	"time"
)

// Нуждается в рефакторинге с использованием дженериков

// FileSet - множество для упрощенных файлов.
//Нужна для того, чтобы определять повторяющиеся файлы
type FileSet struct {
	Files []SimplifiedFile
}

func (set *FileSet) Add(file SimplifiedFile) {
	set.Files = append(set.Files, file)
}

func (set *FileSet) Contains(file SimplifiedFile) bool {
	// Возможно следует использовать что-то вроде cmp.Equal или Reflect.DeepEqual
	for i := range set.Files {
		if set.Files[i].Size == file.Size && set.Files[i].Type == file.Type && set.Files[i].Width == file.Width && set.Files[i].Height == file.Height {
			return true
		}
	}
	return false
}

func (set *FileSet) Find(file SimplifiedFile) *SimplifiedFile {
	// Возможно следует использовать что-то вроде cmp.Equal или Reflect.DeepEqual
	for i := range set.Files {
		if set.Files[i].Size == file.Size && set.Files[i].Type == file.Type && set.Files[i].Width == file.Width && set.Files[i].Height == file.Height {
			return &set.Files[i]
		}
	}
	return nil
}

// FilenameSet - множество имен файлов.
//Нужна для того, чтобы определять повторяющиеся названия файлов
type FilenameSet struct {
	Names []string
}

func (set *FilenameSet) Add(filename string) {
	set.Names = append(set.Names, filename)
}

func (set *FilenameSet) Contains(filename string) bool {
	// Возможно следует использовать что-то вроде cmp.Equal или Reflect.DeepEqual
	for i := range set.Names {
		if set.Names[i] == filename {
			return true
		}
	}
	return false
}

func GenerateFilename() string {
	rand.Seed(time.Now().UnixNano())
	dict := `QWERTYUIOPLKJHGFDSAZXCVBNM1234567890qwertyuioplkjhgfdsazxcvbnm_-`
	var filename string
	for i := 0; i < 10; i++ {
		filename += string(dict[rand.Int()%64])
	}
	return filename + ".jpg"
}

package cmd

import (
	"2arch/internal/parse"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	htmlUrl string
	flags   parse.Flags

	rootCmd = &cobra.Command{
		Use:   "2arch",
		Short: "Архиватор тредов",
		Long: `Архиватор тредов, можно загружать одновременно несколько
Ссылка на тред должна быть вида http(s)://2ch.hk/борда/res/тред.html
Команда для скачивания треда: 2arch -u ссылка_на_тред
Пример: 2arch -u https://2ch.hk/abu/res/42375.html`,
		Run: func(cmd *cobra.Command, args []string) {
			if flags.IsToJson {
				parse.DownloadJson(htmlUrl)
			} else {
				parse.DownloadHtml(htmlUrl, flags)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&htmlUrl, "url", "u", "", "Скачать тред по ссылке")
	rootCmd.PersistentFlags().BoolVarP(&flags.IsToJson, "json", "j", false, "Использовать json")
	rootCmd.PersistentFlags().BoolVarP(&flags.ImagesOnly, "images", "i", false, "Скачивать только картинки")
	rootCmd.PersistentFlags().BoolVarP(&flags.VideosOnly, "videos", "v", false, "Скачивать только видео")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

}

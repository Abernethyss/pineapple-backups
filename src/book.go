package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/book"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	"os"
	"path"
	"strconv"
	"strings"
)

type BookInits struct {
	BookID      string
	ShowBook    bool
	Locks       *threading.GoLimit
	EpubSetting *epub.Epub
}
type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  string
	NovelCover string
	AuthorName string
	CharCount  string
	SignStatus string
}

func (books *BookInits) InitEpubFile() {
	AddImage := true                                                        // add image to epub file
	books.EpubSetting = epub.NewEpub(config.Current.NewBooks["novel_name"]) // set epub setting and add section
	books.EpubSetting.SetAuthor(config.Current.NewBooks["author_name"])     // set author
	if !config.Exist(config.Current.CoverPath) {
		if reader := request.Request(config.Current.NewBooks["novel_cover"]); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(config.Current.CoverPath, "")
		books.EpubSetting.SetCover(strings.ReplaceAll(config.Current.CoverPath, "cover", "../images"), "")
	}

}

func SettingBooks(book_id string) Catalogue {
	var err error
	switch config.Vars.AppType {
	case "sfacg":
		BookInfo := book.GET_BOOK_INFORMATION(book_id)
		if BookInfo.Status.HTTPCode == 200 {
			config.Current.NewBooks = map[string]string{
				"novel_name":  tools.RegexpName(BookInfo.Data.NovelName),
				"novel_id":    strconv.Itoa(BookInfo.Data.NovelID),
				"novel_cover": BookInfo.Data.NovelCover,
				"author_name": BookInfo.Data.AuthorName,
				"char_count":  strconv.Itoa(BookInfo.Data.CharCount),
				"mark_count":  strconv.Itoa(BookInfo.Data.MarkCount),
			}
			err = nil
		} else {
			err = fmt.Errorf(BookInfo.Status.Msg.(string))
		}
	case "cat":
		err = hbooker.GET_BOOK_INFORMATION(book_id)
	}
	if err != nil {
		return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
	}
	fmt.Println(config.Current.NewBooks)
	OutputPath := tools.Mkdir(path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"]))
	config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.Current.NewBooks["novel_name"])
	config.Current.OutputPath = path.Join(OutputPath, config.Current.NewBooks["novel_name"]+".txt")
	config.Current.CoverPath = path.Join("cover", config.Current.NewBooks["novel_name"]+".jpg")
	books := BookInits{BookID: book_id, Locks: nil, ShowBook: true}
	return books.BookDetailed()

}

func (books *BookInits) BookDetailed() Catalogue {
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_id"], config.Current.NewBooks["author_name"],
		config.Current.NewBooks["char_count"],
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	file.Open(config.Current.OutputPath, briefIntroduction, "w")
	return Catalogue{Test: true, EpubSetting: books.EpubSetting}
}

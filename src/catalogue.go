package src

import (
	"fmt"
	"os"
	"sf/src/boluobao"
	"sf/src/config"
	"strconv"
)

type AutoGenerated struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Index    int    `json:"index"`
	IsVip    bool   `json:"is_vip"`
	VolumeID string `json:"volume_id"`
	Content  string `json:"content"`
}

func GetCatalogue(BookData Books) {
	response := boluobao.GetCatalogueDetailedById(BookData.NovelID)
	for _, data := range response.Data.VolumeList {
		fmt.Println("start download volume: ", data.Title)
		for _, Chapter := range data.ChapterList {
			if Chapter.OriginNeedFireMoney == 0 {
				GetContent(len(data.ChapterList), BookData, strconv.Itoa(Chapter.ChapID))
			}
		}
	}
	fmt.Println("NovelName:", BookData.NovelName, "download complete!")
}

func GetContent(ChapterLength int, BookData Books, cid string) {
	response := boluobao.GetContentDetailedByCid(cid)
	if response.Status.HTTPCode != 200 {
		if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" {
			fmt.Println(response.Status.Msg)
			os.Exit(0)
		} else {
			fmt.Println(response.Status.Msg)
		}
	} else {
		if f, err := os.OpenFile(config.Var.SaveFile+"/"+BookData.NovelName+".txt",
			os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			defer func(f *os.File) {
				err = f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(f)
			if _, ok := f.WriteString("\n\n\n" +
				response.Data.Title + ":" + response.Data.AddTime + "\n" +
				response.Data.Expand.Content + "\n" + BookData.AuthorName,
			); ok != nil {
				fmt.Println(ok)
			}
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf(" %d/%d \r", response.Data.ChapOrder, ChapterLength)
}

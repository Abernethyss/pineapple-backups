package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	"strings"
)

type RECOMMEND struct {
	book_list        []string
	recommend_list   [][]string
	book_list_string string
}

func NEW_RECOMMEND() *RECOMMEND {
	var recommend_list [][]string
	recommend := struct {
		Code string `json:"code"`
		Data struct {
			ModuleList []struct {
				ModuleType string `json:"module_type"`
				BossModule struct {
					DesBookList []struct {
						BookID          string `json:"book_id"`
						BookName        string `json:"book_name"`
						CategoryIndex   string `json:"category_index"`
						Description     string `json:"description"`
						AuthorName      string `json:"author_name"`
						Cover           string `json:"cover"`
						DiscountEndTime string `json:"discount_end_time"`
					} `json:"des_book_list"`
				} `json:"boss_module,omitempty"`
			} `json:"module_list"`
		} `json:"data"`
		Tip any `json:"tip"`
	}{}
	request.NewHttpUtils(hbooker.BOOKCITY_RECOMMEND_DATA, "POST").
		Add("theme_type", "NORMAL").Add("tab_type", "200").NewRequests().Unmarshal(&recommend)
	if recommend.Code != "100000" {
		fmt.Println(recommend.Tip.(string))
	} else {
		for _, data := range recommend.Data.ModuleList {
			if data.ModuleType == "1" {
				for _, book := range data.BossModule.DesBookList {
					recommend_list = append(recommend_list, []string{book.BookName, book.BookID})
				}
			}
		}
	}
	return &RECOMMEND{recommend_list: recommend_list}

}

func (is *RECOMMEND) InitBookIdList() {
	is.book_list = nil
	for index, value := range is.recommend_list {
		fmt.Println("index:", index, "\t\tbook id:", value[1], "\t\tbook name:", value[0])
		is.book_list = append(is.book_list, value[1])
	}
	is.book_list_string = strings.Join(is.book_list, ",")
}

func (is *RECOMMEND) CHANGE_NEW_RECOMMEND() {
	change_struct := struct {
		Code string `json:"code"`
		Tip  string `json:"tip"`
		Data struct {
			BookList []struct {
				BookID          string `json:"book_id"`
				BookName        string `json:"book_name"`
				Description     string `json:"description"`
				AuthorName      string `json:"author_name"`
				Cover           string `json:"cover"`
				DiscountEndTime string `json:"discount_end_time"`
				UpStatus        string `json:"up_status"`
				TotalWordCount  string `json:"total_word_count"`
				Introduce       string `json:"introduce"`
			} `json:"book_list"`
		} `json:"data"`
	}{}
	request.NewHttpUtils(hbooker.GET_CHANGE_RECOMMEND, "POST").
		Add("book_id", is.book_list_string).Add("from_module_name", "长篇好书").NewRequests().Unmarshal(&change_struct)
	is.recommend_list = nil
	if change_struct.Code != "100000" {
		fmt.Println(change_struct.Tip)
	} else {
		for _, book := range change_struct.Data.BookList {
			is.recommend_list = append(is.recommend_list, []string{book.BookName, book.BookID})
		}
	}
}

func (is *RECOMMEND) GET_HBOOKER_RECOMMEND() string {
	is.InitBookIdList() // init book_list_string and print recommend_list
	fmt.Println("y is next item recommendation\nd is download recommend book\npress any key to exit..")
	InputChoice := tools.InputStr(">")
	if InputChoice == "y" {
		is.CHANGE_NEW_RECOMMEND() // change recommend_list
		return is.GET_HBOOKER_RECOMMEND()
	} else if InputChoice == "d" {
		return is.book_list[tools.InputInt("input index:", len(is.book_list))]
	}
	return ""

}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"time"
)

type Target struct {
	Url        string `json:"url"`
	Nation     string `json:"nation"`
	CategoryId int    `json:"category Id"`
	Page       int    `json:"page"`
	Per        int    `json:"per"`
	Sort       string `json:"sort"`
}

type Book struct {
	Img          string
	Name         string
	Author       string
	Date         string
	Price        string
	Introduction string
}

var setting = Target{
	Url:        "https://product.kyobobook.co.kr/category",
	Nation:     "KOR",
	CategoryId: 33,
	Page:       1,
	Per:        20,
	Sort:       "new",
}

func Read() {
	// comment slice 생성
	books := make([]*Book, 0)

	// Instantiate default collector
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		RandomDelay: 1 * time.Second,
	})

	// Extract comment|
	//div#homeTab은 정상 작동 div#homeTabNew는 작동 안함 => API 요청이기 때문
	// => 다음 API로 바꾸자 https://product.kyobobook.co.kr/api/gw/pdt/category/new?page=1&per=20&saleCmdtDvsnCode=KOR&saleCmdtClstCode=33&sort=new
	c.OnHTML("div#homeTabNew .prod_list li.prod_item", func(e *colly.HTMLElement) {

		b := &Book{}

		b.Img = e.ChildAttr("span.img_box img", "src")
		b.Name = e.ChildText(".prod_name")
		b.Author = e.ChildText(".prod_author")
		b.Date = e.ChildText("span.date")
		b.Price = e.ChildText(".prod_price .val")
		b.Introduction = e.ChildText(".prod_introduction")

		fmt.Println("element unmarshal")
		fmt.Println(b)

		books = append(books, b)

	})

	// 핸들 다 걸었으면 요청
	// 교보문고에서
	var url = fmt.Sprintf("%s/%s/%d#?page=%d&per=%d&sort=%s", setting.Url, setting.Nation, setting.CategoryId, setting.Page, setting.Per, setting.Sort)
	c.Visit(url)

	// json 관련 코드
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(books)
	fmt.Println("url :", url)
}

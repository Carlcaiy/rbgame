package main

import (
	"fmt"

	"log"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	// 请求html页面

	res, err := http.Get("https://www.emojiall.com/zh-hant/emoji-art-list?page=1")

	if err != nil {

		// 错误处理

		log.Fatal(err)

	}

	defer res.Body.Close()

	if res.StatusCode != 200 {

		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)

	}

	// 加载 HTML document对象

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {

		log.Fatal(err)

	}

	// Find the review items
	fmt.Println("find")
	doc.Find(".w-100 .h-100 .fontsize_1x").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Add(".col-auto").Add("a").Text())
	})

}

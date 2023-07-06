package main

import (
	"fmt"

	"log"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	// 请求html页面

	res, err := http.Get("http://metalsucks.net")

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

	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {

		// For each item found, get the band and title

		band := s.Find("a").Text()

		title := s.Find("i").Text()

		fmt.Printf("Review %d: %s - %s\n", i, band, title)

	})

}

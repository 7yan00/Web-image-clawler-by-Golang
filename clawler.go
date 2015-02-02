package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"flag"
	"github.com/PuerkitoBio/goquery"
)

func GetPage(base string) []*url.URL {
	doc, _ := goquery.NewDocument(base)
	var result []*url.URL
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		target, _ := s.Attr("src")
		base,_ := url.Parse(base)
		targets, _ := url.Parse(target)
		result = append(result, base.ResolveReference(targets))
	})
	return result
}

func GetImage(urls []*url.URL) {
	for i, url := range urls {
		urrl := url.String()
		response, err := http.Get(urrl)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		file, err := os.Create(fmt.Sprintf("hoge%d.jpg", i))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		io.Copy(file, response.Body)
	}
}

func main() {
	flag.Parse()
	fmt.Println("It works!")
	base := flag.Arg(0)
	urls := (GetPage(base))
	GetImage(urls)
	fmt.Println("fin! You ought to get some picture!")
}

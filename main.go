package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var stock = []string{}
var base = " "
var i int = 0
var wg = new(sync.WaitGroup)

func main() {
	flag.Parse()
	fmt.Println("It works!")
	base = flag.Arg(0)
	doc, _ := goquery.NewDocument(base)
	results := makeUrl(doc)
	for len(results) > 0 {
		results = GetUrl(results)
	}
	wg.Wait()
}

func containsInStock(value string) bool {
	l := len(stock)
	for i := 0; i < l; i++ {
		if stock[i] == value {
			return true
		}
	}
	return false
}

func GetUrl(urls []*url.URL) []*url.URL {
	sarceUrl := []*url.URL{}
L:
	for _, item := range urls {
		url_string := item.String()
		if !strings.Contains(url_string, base) {
			continue L
		}
		if containsInStock(url_string) {
			continue L
		}
		fmt.Println(url_string)
		stock = append(stock, url_string)
		doc, _ := goquery.NewDocument(url_string)
		results := makeUrl(doc)
		wg.Add(1)
		go GetImage(doc)
		sarceUrl = append(sarceUrl, results...)
	}
	return sarceUrl
}

func makeUrl(doc *goquery.Document) []*url.URL {
	var result []*url.URL
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		target, _ := s.Attr("href")
		base, _ := url.Parse(base)
		targets, _ := url.Parse(target)
		result = append(result, base.ResolveReference(targets))
	})
	return result
}

func GetImage(doc *goquery.Document) {
	var result []*url.URL
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		target, _ := s.Attr("src")
		base, _ := url.Parse(base)
		targets, _ := url.Parse(target)
		result = append(result, base.ResolveReference(targets))
	})
	for _, imageUrl := range result {
		imageUrl_String := imageUrl.String()
		if containsInStock(imageUrl_String) {
			continue
		}
		response, err := http.Get(imageUrl_String)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		file, err := os.Create(fmt.Sprintf("hoge%d.jpg", i))
		i++
		if err != nil {
			panic(err)
		}
		defer file.Close()
		io.Copy(file, response.Body)
	}
	wg.Done()
}

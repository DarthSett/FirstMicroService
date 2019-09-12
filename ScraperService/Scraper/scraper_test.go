package Scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
	"time"
)

func TestScrape(t *testing.T) {
	link := "https://www.amazon.in/Korecall-KORUSB2-Phone-Recorder-Black/dp/B00N5TL0T8/ref=sr_1_1?qid=1568273090&s=electronics&sr=1-1"
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Timeout: 30 * time.Second}
	req.Header.Set("User-Agent", "Not Firefox")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	t.Run("It scrapes the title properly", func(t *testing.T) {
		title := titleScraper(doc)
		AssertEqual(t, "Korecall KORUSB2 2 line USB Phone Recorder (Black or White)", title, "")
	})
	t.Run("It scrapes the price properly", func(t *testing.T) {
		price := priceScraper(doc)
		AssertEqual(t, "4800.00", price, "")
	})
	t.Run("It scrapes the title properly", func(t *testing.T) {
		seller := sellerScraper(doc)
		AssertEqual(t, "Realtime Solutions", seller, "")
	})
}
func AssertEqual(t *testing.T, expected interface{}, got interface{}, message string) {
	if expected == got {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", expected, got)
	}
	t.Fatal(message)
}

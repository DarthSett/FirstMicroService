package Scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"testing"
	"time"
)

func TestScrape(t *testing.T) {
	link := "https://www.amazon.in/Power-Mens-Lionel-Running-Shoes/dp/B01FQXR2ME?pf_rd_p=4560c9eb-731c-49c7-bd28-5432424e2e3c&pd_rd_wg=65Spy&pf_rd_r=1NE5C8E8FHAXWP9B5J6C&ref_=pd_gw_unk&pd_rd_w=KFmjp&pd_rd_r=9644ad67-0f90-4ba4-9eb1-d17f2188c441"
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{Timeout: 30 * time.Second}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.47 Safari/537.36")
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
		AssertEqual(t, "Power Men's Lionel Running Shoes", title, "")
	})
	t.Run("It scrapes the price properly", func(t *testing.T) {
		price := priceScraper(doc)
		AssertEqual(t, "1099.00", price, "")
	})
	t.Run("It scrapes the title properly", func(t *testing.T) {
		seller := sellerScraper(doc)
		AssertEqual(t, "Craftnation", seller, "")
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

package Scraper

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Scrape(db *sql.DB) {

	println("In Scraper")
	row := dbGetLink(db)
	defer row.Close()
	println("links fetched")
	client := http.Client{Timeout: 30 * time.Second}
	link := ""
	status := ""
	println("Entering loop")
	for row.Next() {

		err := row.Scan(&link, &status)
		FailOnError(err, "Failed to scan the fetched row")

		//res, err := http.Get(link)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//println("\n\n\n")
		//b ,_ := ioutil.ReadAll(res.Body)
		//println(string(b))
		//println("\n\n\n")
		println("updating link to scraping")
		db.Query("update Link set status = 'scraping' where link = ?", link)
		req, err := http.NewRequest(http.MethodGet, link, nil)
		FailOnError(err, "failed to create a req")
		req.Header.Set("User-Agent", "Not Firefox")
		println("1")
		res, err := client.Do(req)
		FailOnError(err, "Failed to send a request")
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		println("2")
		doc, err := goquery.NewDocumentFromReader(res.Body)
		FailOnError(err, "Failed to read file from body")
		e := "error while scraping: "
		flag := false
		title := titleScraper(doc)
		println("3")
		if title == "" {
			println("#######")
			flag = true
			e = e + "tag 'span#productTitle' not found, "
			errScrape := fmt.Sprintf(`update Link set status = "%s" where link = "%s"`, e, link)
			_, err = db.Query(errScrape)
			FailOnError(err, "failed to update db")
		}
		p := priceScraper(doc)
		price, _ := strconv.ParseFloat(p, 64)
		println("4")
		if p == "" {
			println("!!!!!!")
			flag = true
			e = e + "tag 'data-asin-price' not found, "
			errScrape := fmt.Sprintf(`update Link set status = "%s" where link = "%s"`, e, link)
			_, err = db.Query(errScrape)
			FailOnError(err, "failed to update db")
		}
		seller := sellerScraper(doc)
		println("5")
		if seller == "" {
			println("@@@@@@")
			flag = true
			e = e + "tag 'sellerProfileTriggerId' not found"
			errScrape := fmt.Sprintf(`update Link set status = "%s" where link = "%s"`, e, link)
			_, err = db.Query(errScrape)
			FailOnError(err, "failed to update db")
		}
		fmt.Printf("\ntitle: %s \nprice: %.2f \nseller: %s", title, price, seller)
		dbInsert(db, title, price, seller)
		if !flag {
			db.Query("update Link set status = 'scraped' where link = ?", link)
		}
		res.Body.Close()
	}

}

func dbInsert(db *sql.DB, title string, price float64, seller string) {

	println("6")
	q := fmt.Sprintf(`INSERT into Products values ("%s", "%.2f", "%s")`, title, price, seller)
	_, err := db.Query(q)
	FailOnError(err, "Failed to connect to mysql")

}

func dbGetLink(db *sql.DB) *sql.Rows {
	println("getting links")
	row, err := db.Query("select * from Link where status = 'unscraped'")
	FailOnError(err, "Failed to fetch links from db")
	return row
}

func titleScraper(doc *goquery.Document) string {
	title := doc.Find("span#productTitle").Text()
	title = string(bytes.TrimSpace([]byte(title)))
	return title
}

func sellerScraper(doc *goquery.Document) string {
	seller := doc.Find("a#sellerProfileTriggerId").Text()
	seller = string(bytes.TrimSpace([]byte(seller)))
	return seller
}

func priceScraper(doc *goquery.Document) string {
	p, _ := doc.Find("div#cerberus-data-metrics").Attr("data-asin-price")
	p = string(bytes.TrimSpace([]byte(p)))
	return p
}

func FailOnError(err error, msg string) {
	if err != nil {
		println(msg + ": " + err.Error())
	}
}

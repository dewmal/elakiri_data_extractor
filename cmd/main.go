package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/url"
	"time"
	"webcrawler/cmd/data"
	"webcrawler/cmd/extractor"
)

func main() {

	dsn := url.URL{
		User:     url.UserPassword("postgres", "dewmal91"),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", "localhost", 5432),
		Path:     "ek_crawler_db",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	db, err := gorm.Open("postgres", dsn.String())
	db.LogMode(true)
	if err != nil {
		println(err)
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&data.UserPost{})
	db.AutoMigrate(&data.UserProfile{})

	c := colly.NewCollector(
		colly.AllowedDomains("elakiri.com", "www.elakiri.com"),
		colly.CacheDir("./ek_cache"),
	)
	c.Limit(
		&colly.LimitRule{
			DomainGlob:  "*elakiri.*",
			RandomDelay: 2 * time.Second,
			Parallelism: 10,
		})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		uniqueId, _ := uuid.Parse(e.Request.URL.String())
		fmt.Println(uniqueId)
	})

	postDataCollector := c.Clone()
	userDataCollector := c.Clone()

	visitURL := func(pageUrlString string) {
		pageUrl, _ := url.Parse(pageUrlString)
		if pageUrl.Path == "/forum/showthread.php" {
			postDataCollector.Visit(pageUrlString)
		} else if pageUrl.Path == "/forum/member.php" {
			userDataCollector.Visit(pageUrlString)
		} else {
			//c.Visit(pageUrlString)
		}
	}

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		var pageUrlString = e.Request.AbsoluteURL(e.Attr("href"))
		visitURL(pageUrlString)
	})
	postDataCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		var pageUrlString = e.Request.AbsoluteURL(e.Attr("href"))
		visitURL(pageUrlString)
	})
	userDataCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		var pageUrlString = e.Request.AbsoluteURL(e.Attr("href"))
		visitURL(pageUrlString)
	})

	userDataCollector.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println("Extracting user Data ", e.Request.URL.Query().Get("u"))
		extractor.ExtractUserDetails(e, db)
	})
	postDataCollector.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println("Extracting Thread Data ", e.Request.URL.Query().Get("t"))
		extractor.ExtractThreadDetail(e, db)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	postDataCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Thread", r.URL)
	})
	userDataCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting User Profile", r.URL)
	})

	c.Visit("http://www.elakiri.com")

	c.Wait()
}

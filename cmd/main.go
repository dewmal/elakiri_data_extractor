package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/url"
	"time"
)

func olmain() {
	c := colly.NewCollector(
		colly.AllowedDomains("elakiri.com", "www.elakiri.com"),
		colly.CacheDir("./ek_cache"),
	)
	c.Limit(
		&colly.LimitRule{
			DomainGlob:  "*elakiri.*",
			RandomDelay: 10 * time.Second,
			Parallelism: 2,
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

	c.Visit("http://www.elakiri.com/forum/showthread.php?t=1937671")

	c.Wait()
}

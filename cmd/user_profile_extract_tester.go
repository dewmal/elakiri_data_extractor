package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/url"
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

	c := colly.NewCollector()
	c.OnHTML("body", func(be *colly.HTMLElement) {
		extractor.ExtractUserDetails(be, db)
	})
	c.Visit("http://www.elakiri.com/forum/member.php?s=6c6b0f832ae57f4674e3cb8384e94947&u=189162")
}

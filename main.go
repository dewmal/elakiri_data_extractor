package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"log"
	"net/url"
	"os"
	"time"
	"webcrawler/cmd/data"
	"webcrawler/cmd/extractor"
)

type Config struct {
	Crawler struct {
		DB struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			UserName string `yaml:"username"`
			Password string `yaml:"password"`
			DataBase string `yaml:"dataBase"`
		}
	}
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", os.Getenv("CRAWLER_CONFIG_PATH"), "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func main() {
	// by the user in the flags
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	dsn := url.URL{
		User:     url.UserPassword(cfg.Crawler.DB.UserName, cfg.Crawler.DB.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", cfg.Crawler.DB.Host, cfg.Crawler.DB.Port),
		Path:     cfg.Crawler.DB.DataBase,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	log.Println(dsn.String())
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
	db.AutoMigrate(&data.Thread{})

	c := colly.NewCollector(
		colly.AllowedDomains("elakiri.com", "www.elakiri.com"),
		colly.CacheDir("./ek_cache"),
	)
	c.Limit(
		&colly.LimitRule{
			DomainGlob:  "*elakiri.*",
			RandomDelay: 4 * time.Second,
			Parallelism: 20,
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

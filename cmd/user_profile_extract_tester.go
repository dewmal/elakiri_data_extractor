package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/url"
	"strconv"
	"time"
	"webcrawler/cmd/data"
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

	c := colly.NewCollector()
	c.OnHTML("body", func(be *colly.HTMLElement) {

		var statList []string
		be.ForEach(".profilefield_list dd", func(i int, element *colly.HTMLElement) {
			statList = append(statList, element.Text)
		})
		joinDateVal := statList[0]
		totalPost, _ := strconv.ParseInt(statList[1], 0, 0)
		joinDate, _ := time.Parse("01-02-2006", joinDateVal)

		userName := be.ChildText("#username_box h1")

		userProfile := data.UserProfile{
			UserName:    userName,
			UserId:      "",
			JoinDateVal: joinDateVal,
			JoinDate:    joinDate,
			TotalPost:   totalPost,
		}
		fmt.Println(userProfile)
		// User Visitor Messages
		be.ForEach("html body table tbody tr td div div.page div div#usercss.floatcontainer div#content_container div#content div#profile_tabs div#visitor_messaging.tborder.content_block div#collapseobj_visitor_messaging.block_content ol#message_list.alt1.block_row.list_no_decoration li", func(i int, element *colly.HTMLElement) {
			usernameText := element.ChildText("a.username")
			messageBody := element.ChildText(".visitor_message_body")
			messageTimeString := element.ChildText(".visitor_message_date")
			messageTime, _ := time.Parse("01-02-2006 03:04 PM", messageTimeString) //11-28-2019 11:31 AM
			userLink, _ := url.Parse(element.ChildAttr("a.username", "href"))

			userId := userLink.Query().Get("u")
			if userId != "" {
				vm := data.UserPost{
					Username:      usernameText,
					UserId:        userId,
					Message:       messageBody,
					MessageSource: messageBody,
					PostTimeVal:   messageTimeString,
					PostTime:      messageTime,
					PostType:      data.PostTypeEnum.VisitorPost,
				}
				println(vm.Username)
				//db.Create(&vm)
			}

		})
	})
	c.Visit("http://www.elakiri.com/forum/member.php?s=6c6b0f832ae57f4674e3cb8384e94947&u=189162")
}

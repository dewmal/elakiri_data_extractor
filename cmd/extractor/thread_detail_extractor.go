package extractor

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"webcrawler/cmd/data"
)

/**
Extract And Store Profile Details
*/
func ExtractThreadDetail(be *colly.HTMLElement, db *gorm.DB) {
	//db := dbVal.Begin()
	pageUrl := be.Request.URL
	if pageUrl.Path != "/forum/showthread.php" {
		return //errors.New("invalid url")
	}
	baseURL := be.Request.URL
	pageId := baseURL.Query().Get("page")
	isOpenPage := pageId == "" || pageId == "1"
	threadId, _ := strconv.ParseInt(baseURL.Query().Get("t"), 0, 0)

	var ownerName string
	var ownerId int64
	var postTime time.Time

	be.ForEachWithBreak("div#posts div div.page", func(i int, element *colly.HTMLElement) bool {

		usernameText := element.ChildText("a.bigusername")
		userLink, _ := url.Parse(element.ChildAttr("a.bigusername", "href"))
		userId, _ := strconv.ParseInt(userLink.Query().Get("u"), 0, 0)

		userLocation := strings.Replace(element.ChildTexts("div.postbit_box div")[2], "Location:", "", -1)

		messageBody := element.ChildTexts("div.vb_postbit")[0]

		//relatedPosts := element.ChildTexts("vb_postbit")
		var relatedPosts []int64

		element.ForEach("div.vb_postbit", func(i int, element *colly.HTMLElement) {
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				quoteUrl, _ := url.Parse(element.Attr("href"))
				//fmt.Println(quoteUrl.String())
				if strings.HasSuffix(quoteUrl.Path, "showthread.php") {
					//fmt.Println("Show Thread ",quoteUrl.String())
					quoteId, _ := strconv.ParseInt(quoteUrl.Query().Get("p"), 10, 64)
					relatedPosts = append(relatedPosts, quoteId)
				}
			})

		})

		var messageBodySource string
		element.ForEach("div.vb_postbit", func(i int, element *colly.HTMLElement) {
			htmlVal, _ := element.DOM.Html()
			messageBodySource += htmlVal
		})
		// Extract Quoted Texts

		// Extract Time
		messageTimeStringReplaceVal := element.ChildText("div div div table.tborder tbody tr td.alt1 div.smallfont strong")
		messageTimeStringRawVal := element.ChildTexts("div div div table.tborder tbody tr td.alt1 div.smallfont")[0]
		messageTimeString := strings.TrimSpace(strings.Replace(strings.Replace(messageTimeStringRawVal, messageTimeStringReplaceVal, "", -1), ",", "", -1))

		messageTime, _ := time.Parse("01-02-2006 03:04 PM", messageTimeString) //11-28-2019 11:31 AM

		postLink, _ := url.Parse(element.ChildAttr("div div table.tborder tbody tr td.thead div.smallfont a", "href"))
		postId, _ := strconv.ParseInt(postLink.Query().Get("p"), 0, 0)
		if len(relatedPosts) > 0 {
			log.Println("Related Post with  ", postId)
		}

		var up data.UserPost
		db.Where(&data.UserPost{
			PostId: postId,
		}).FirstOrCreate(&up)

		up.ThreadId = threadId
		up.PostId = postId
		up.Username = usernameText
		up.UserId = userId
		up.Message = messageBody
		up.MessageSource = messageBodySource
		up.PostTimeVal = messageTimeString
		up.PostTime = messageTime
		up.PostType = data.PostTypeEnum.ThreadPost
		up.ThreadId = threadId
		up.RelatedPosts = relatedPosts

		db.Save(&up)

		var user data.UserProfile
		db.Where(&data.UserProfile{UserId: userId}).FirstOrCreate(&user)
		user.Location = userLocation
		user.UserName = usernameText
		db.Save(&user)

		if isOpenPage && i == 0 {
			ownerName = user.UserName
			ownerId = user.UserId
			postTime = messageTime
		}
		return true
	})

	var tags []string
	be.ForEach("td#tag_list_cell a", func(i int, element *colly.HTMLElement) {
		tags = append(tags, element.Text)
	})

	var categories []string
	be.ForEach("div.page div table.tborder tbody tr td.alt1 table tbody tr td span.smallfont a", func(i int, element *colly.HTMLElement) {
		categories = append(categories, element.Text)
	})

	title := be.ChildText("div.page div table.tborder tbody tr td.alt1 table tbody tr td.smallfont strong")

	var thread data.Thread
	db.Where(&data.Thread{
		ThreadId: threadId,
	}).FirstOrCreate(&thread)
	thread.Tags = tags
	thread.OwnerId = ownerId
	thread.OwnerUser = ownerName
	thread.Categories = categories
	thread.Title = title

	if postTime != (time.Time{}) {
		thread.PostDateTime = postTime
	}

	db.Save(&thread)
	fmt.Println(thread.Title)
	//return db.Commit().Error
}

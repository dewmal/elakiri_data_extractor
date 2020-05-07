package extractor

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/jinzhu/gorm"
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

	baseURL := be.Request.URL
	pageId := baseURL.Query().Get("page")
	isOpenPage := pageId == "" || pageId == "1"
	threadId, _ := strconv.ParseInt(baseURL.Query().Get("t"), 0, 0)

	var ownerName string
	var ownerId int64
	var postTime time.Time

	be.ForEachWithBreak("#posts div div.page", func(i int, element *colly.HTMLElement) bool {

		usernameText := element.ChildText("a.bigusername")
		userLink, _ := url.Parse(element.ChildAttr("a.bigusername", "href"))
		userId, _ := strconv.ParseInt(userLink.Query().Get("u"), 0, 0)

		userLocation := strings.Replace(element.ChildTexts("div.postbit_box div")[2], "Location:", "", -1)

		messageBody := "" //element.ChildTexts("div.vb_postbit")
		// Extract Quoted Texts

		// Extract Time
		messageTimeStringReplaceVal := element.ChildText("div div div table.tborder tbody tr td.alt1 div.smallfont strong")
		messageTimeStringRawVal := element.ChildTexts("div div div table.tborder tbody tr td.alt1 div.smallfont")[0]
		messageTimeString := strings.TrimSpace(strings.Replace(strings.Replace(messageTimeStringRawVal, messageTimeStringReplaceVal, "", -1), ",", "", -1))

		messageTime, _ := time.Parse("01-02-2006 03:04 PM", messageTimeString) //11-28-2019 11:31 AM

		postLink, _ := url.Parse(element.ChildAttr("td.thread div a", "href"))
		postId := postLink.Query().Get("p")

		up := data.UserPost{
			PostId:        postId,
			Username:      usernameText,
			UserId:        userId,
			Message:       messageBody,
			MessageSource: messageBody,
			PostTimeVal:   messageTimeString,
			PostTime:      messageTime,
			PostType:      data.PostTypeEnum.ThreadPost,
			ThreadId:      threadId,
		}

		db.Save(&up)

		var user data.UserProfile
		db.Where(&data.UserProfile{UserId: userId}).FirstOrInit(&user)
		user.Location = userLocation
		user.UserName = usernameText
		db.Save(&user)

		if isOpenPage && i == 0 {
			ownerName = user.UserName
			ownerId = user.UserId
			postTime = messageTime
		}
		return false
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
	}).FirstOrInit(&thread)
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

}

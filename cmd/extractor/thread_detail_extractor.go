package extractor

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"net/url"
	"strconv"
	"strings"
	"time"
	"webcrawler/cmd/data"
)

/**
Extract And Store Profile Details
*/
func ExtractThreadDetail(be *colly.HTMLElement) (data.Thread, []data.UserPost, []data.UserProfile, error) {
	var dataUserPosts []data.UserPost
	var dataUserProfiles []data.UserProfile
	var dataThread data.Thread = data.Thread{}
	//db := dbVal.Begin()
	pageUrl := be.Request.URL
	if pageUrl.Path != "/forum/showthread.php" {
		return dataThread, nil, nil, errors.New("invalid URL") //errors.New("invalid url")
	}
	baseURL := be.Request.URL
	pageId := baseURL.Query().Get("page")
	isOpenPage := pageId == "" || pageId == "1"
	threadId, _ := strconv.ParseInt(baseURL.Query().Get("t"), 0, 0)

	var ownerName string
	var ownerId int64
	var postTime time.Time

	be.ForEach("div#posts div div.page", func(i int, element *colly.HTMLElement) {

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
				if quoteUrl != nil && strings.HasSuffix(quoteUrl.Path, "showthread.php") {
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

		var dataUserPost data.UserPost = data.UserPost{}
		dataUserPost.ThreadId = threadId
		dataUserPost.PostId = postId
		dataUserPost.Username = usernameText
		dataUserPost.UserId = userId
		dataUserPost.Message = messageBody
		dataUserPost.MessageSource = messageBodySource
		dataUserPost.PostTimeVal = messageTimeString
		dataUserPost.PostTime = messageTime
		dataUserPost.PostType = data.PostTypeEnum.ThreadPost
		dataUserPost.ThreadId = threadId
		dataUserPost.RelatedPosts = relatedPosts

		dataUserPosts = append(dataUserPosts, dataUserPost)

		var user data.UserProfile = data.UserProfile{}
		user.Location = userLocation
		user.UserName = usernameText
		user.UserId = userId

		dataUserProfiles = append(dataUserProfiles, user)

		if isOpenPage && i == 0 {
			ownerName = user.UserName
			ownerId = user.UserId
			postTime = messageTime
		}
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

	dataThread.ThreadId = threadId
	dataThread.Tags = tags
	dataThread.OwnerId = ownerId
	dataThread.OwnerUser = ownerName
	dataThread.Categories = categories
	dataThread.Title = title

	if postTime == (time.Time{}) {
		dataThread.PostDateTime = postTime
	}

	return dataThread, dataUserPosts, dataUserProfiles, nil
}

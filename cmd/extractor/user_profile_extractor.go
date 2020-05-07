package extractor

import (
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
func ExtractUserDetails(be *colly.HTMLElement, db *gorm.DB) {
	userProfileId := be.Request.URL.Query().Get("u")

	var statList []string
	be.ForEach(".profilefield_list dd", func(i int, element *colly.HTMLElement) {
		statList = append(statList, element.Text)
	})
	joinDateVal := statList[0]
	totalPost, _ := strconv.ParseInt(statList[1], 0, 0)
	joinDate, _ := time.Parse("01-02-2006", joinDateVal)

	userName := be.ChildText("#username_box h1")
	memberStatus := be.ChildText("#username_box h2")

	// Reputation Rank
	var posRepCount int
	var negRepCount int

	be.ForEach("#reputation img", func(i int, element *colly.HTMLElement) {
		isPosRep := strings.HasPrefix(element.Attr("src"), "pos.gif")
		if isPosRep {
			negRepCount += 1
		} else {
			posRepCount += 1
		}
	})
	reputationRank := posRepCount - negRepCount

	// Friend Detail
	var friendList []string
	be.ForEach("#friends_list_big li.memberinfo_small", func(i int, ef *colly.HTMLElement) {
		firendUserName := ef.ChildText("a.bigusername")
		firendLink, _ := url.Parse(ef.ChildAttr("a.bigusername", "href"))
		friendId := firendLink.Query().Get("u")

		friendUser := data.UserProfile{
			UserName: firendUserName,
			UserId:   friendId,
		}
		db.Create(&friendUser)
		friendList = append(friendList, friendId)
	})
	// Visitor Detail
	var visitorList []string
	be.ForEach(".last_visitors_list li.smallfont", func(i int, ef *colly.HTMLElement) {

		visitorLink, _ := url.Parse(ef.ChildAttr("a", "href"))
		visitorId := visitorLink.Query().Get("u")

		visitorList = append(visitorList, visitorId)
	})

	// Save user profile
	userProfile := data.UserProfile{
		UserName:       userName,
		UserId:         userProfileId,
		JoinDateVal:    joinDateVal,
		JoinDate:       joinDate,
		TotalPost:      totalPost,
		MemberStatus:   memberStatus,
		ReputationRank: reputationRank,
		Friends:        friendList,
		LastVisitors:   visitorList,
	}
	db.Save(&userProfile)

	// User Visitor Messages
	be.ForEach("html body table tbody tr td div div.page div div#usercss.floatcontainer div#content_container div#content div#profile_tabs div#visitor_messaging.tborder.content_block div#collapseobj_visitor_messaging.block_content ol#message_list.alt1.block_row.list_no_decoration li", func(i int, element *colly.HTMLElement) {
		usernameText := element.ChildText("a.username")
		messageBody := element.ChildText(".visitor_message_body")
		messageTimeString := element.ChildText(".visitor_message_date")
		messageTime, _ := time.Parse("01-02-2006 03:04 PM", messageTimeString) //11-28-2019 11:31 AM
		userLink, _ := url.Parse(element.ChildAttr("a.username", "href"))
		userId := userLink.Query().Get("u")
		postLink, _ := url.Parse(element.ChildAttr("ul li a", "href"))
		postId := postLink.Query().Get("u1")

		if userId != "" {
			vm := data.UserPost{
				Username:      usernameText,
				UserId:        userId,
				Message:       messageBody,
				MessageSource: messageBody,
				PostTimeVal:   messageTimeString,
				PostTime:      messageTime,
				PostType:      data.PostTypeEnum.VisitorPost,
				PostId:        postId,
			}
			db.Create(&vm)
		}
	})
}

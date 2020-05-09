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
	//db := dbVal.Begin()
	pageUrl := be.Request.URL
	if pageUrl.Path != "/forum/member.php" {
		return
		//errors.New("invalid url")
	}
	userProfileId, _ := strconv.ParseInt(be.Request.URL.Query().Get("u"), 0, 0)

	var statList []string
	be.ForEach(".profilefield_list dd", func(i int, element *colly.HTMLElement) {
		statList = append(statList, element.Text)
	})
	var joinDateIndex int
	var TotalPostIndex int
	be.ForEachWithBreak(".profilefield_list dt.shade", func(i int, element *colly.HTMLElement) bool {
		joinDateIndex = i
		if element.Text == "Join Date" {
			return false
		}
		return true
	})
	be.ForEachWithBreak(".profilefield_list dt.shade", func(i int, element *colly.HTMLElement) bool {
		TotalPostIndex = i
		if element.Text == "Total Post" {
			return false
		}
		return true
	})
	var totalPost int64
	var joinDate time.Time
	var joinDateVal string
	if len(statList) > joinDateIndex {
		joinDateVal := statList[joinDateIndex]
		joinDate, _ = time.Parse("01-02-2006", joinDateVal)
	}

	if len(statList) > TotalPostIndex {
		totalPost, _ = strconv.ParseInt(strings.Replace(statList[TotalPostIndex], ",", "", -1), 0, 0)
	}

	userName := be.ChildText("#username_box h1")
	memberStatus := be.ChildText("#username_box h2")

	// Reputation Rank
	var posRepCount int
	var negRepCount int

	be.ForEach("div#usercss.floatcontainer div#content_container div#content div#main_userinfo.floatcontainer table tbody tr td#username_box div#reputation_rank div#reputation img.inlineimg", func(i int, element *colly.HTMLElement) {
		isPosRep := strings.HasSuffix(element.Attr("src"), "pos.gif")
		if isPosRep {
			posRepCount += 1
		} else {
			negRepCount += 1
		}
		if strings.HasSuffix(element.Attr("src"), "balance.gif") {
			posRepCount = 0
			negRepCount = 0
		}
	})
	reputationRank := posRepCount - negRepCount
	//println(reputationRank, posRepCount, negRepCount)

	// Friend Detail
	var friendList []int64
	be.ForEach("#friends_list_big li.memberinfo_small", func(i int, ef *colly.HTMLElement) {
		friendUserName := ef.ChildText("a.bigusername")
		friendLink, _ := url.Parse(ef.ChildAttr("a.bigusername", "href"))
		friendId, _ := strconv.ParseInt(friendLink.Query().Get("u"), 0, 0)

		var friendUser data.UserProfile
		db.Where(&data.UserProfile{
			UserId: friendId,
		}).FirstOrCreate(&friendUser)
		friendUser.UserId = friendId
		friendUser.UserName = friendUserName
		db.Save(&friendUser)
		friendList = append(friendList, friendId)
	})
	// Visitor Detail
	var visitorList []int64
	be.ForEach(".last_visitors_list li.smallfont", func(i int, ef *colly.HTMLElement) {
		visitorLink, _ := url.Parse(ef.ChildAttr("a", "href"))
		visitorId, _ := strconv.ParseInt(visitorLink.Query().Get("u"), 0, 0)
		visitorList = append(visitorList, visitorId)
	})
	// Get Visitor count
	totalVisit := be.ChildText("#collapseobj_visitors div.block_row.block_footer strong")
	totalVisitCount, _ := strconv.ParseInt(strings.Replace(totalVisit, ",", "", -1), 0, 0)
	// Save user profile
	var userProfile data.UserProfile
	db.Where(&data.UserProfile{
		UserId: userProfileId,
	}).FirstOrCreate(&userProfile)
	userProfile.UserName = userName
	userProfile.UserId = userProfileId
	userProfile.JoinDateVal = joinDateVal
	userProfile.JoinDate = joinDate
	userProfile.TotalPost = totalPost
	userProfile.MemberStatus = memberStatus
	userProfile.ReputationRank = reputationRank
	userProfile.Friends = friendList
	userProfile.LastVisitors = visitorList
	userProfile.TotalPageVisit = totalVisitCount

	db.Save(&userProfile)
	//println(userProfile.ReputationRank)

	// User Visitor Messages
	be.ForEach("html body table tbody tr td div div.page div div#usercss.floatcontainer div#content_container div#content div#profile_tabs div#visitor_messaging.tborder.content_block div#collapseobj_visitor_messaging.block_content ol#message_list.alt1.block_row.list_no_decoration li", func(i int, element *colly.HTMLElement) {
		usernameText := element.ChildText("a.username")
		messageBody := element.ChildText(".visitor_message_body")

		var messageBodySource string
		element.ForEach(".visitor_message_body", func(i int, element *colly.HTMLElement) {
			htmlVal, _ := element.DOM.Html()
			messageBodySource += htmlVal
		})

		messageTimeString := element.ChildText(".visitor_message_date")
		messageTime, _ := time.Parse("01-02-2006 03:04 PM", messageTimeString) //11-28-2019 11:31 AM
		userLink, _ := url.Parse(element.ChildAttr("a.username", "href"))
		userId, _ := strconv.ParseInt(userLink.Query().Get("u"), 0, 0)
		postLink, _ := url.Parse(element.ChildAttr("ul li a", "href"))
		postId, _ := strconv.ParseInt(postLink.Query().Get("u1"), 0, 0)

		if userId != 0 {
			var up data.UserPost
			db.Where(&data.UserPost{
				PostId: postId,
			}).FirstOrCreate(&up)

			up.Username = usernameText
			up.UserId = userId
			up.Message = messageBody
			up.MessageSource = messageBodySource
			up.PostTimeVal = messageTimeString
			up.PostTime = messageTime
			up.PostType = data.PostTypeEnum.VisitorPost
			db.Save(&up)
		}
	})
	//return db.Commit().Error
}

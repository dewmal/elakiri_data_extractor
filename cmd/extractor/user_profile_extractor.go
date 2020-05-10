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
func ExtractUserDetails(be *colly.HTMLElement) (data.UserProfile, []data.VisitorMessage, []data.UserProfile, error) {

	var dataFriendList []data.UserProfile
	var dataVisitorPosts []data.VisitorMessage
	var dataUserProfile = data.UserProfile{}

	pageUrl := be.Request.URL
	if pageUrl.Path != "/forum/member.php" {
		return dataUserProfile, nil, nil, errors.New("invalid URL")
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

		friendUser := data.UserProfile{
			UserId:   friendId,
			UserName: friendUserName,
		}
		dataFriendList = append(dataFriendList, friendUser)
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

	dataUserProfile.UserName = userName
	dataUserProfile.UserId = userProfileId
	dataUserProfile.JoinDateVal = joinDateVal
	dataUserProfile.JoinDate = joinDate
	dataUserProfile.TotalPost = totalPost
	dataUserProfile.MemberStatus = memberStatus
	dataUserProfile.ReputationRank = reputationRank
	dataUserProfile.Friends = friendList
	dataUserProfile.LastVisitors = visitorList
	dataUserProfile.TotalPageVisit = totalVisitCount

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

		//println("Post ID ", postId, postLink.String())
		if userId != 0 {

			var up = data.VisitorMessage{}
			up.PostUserName = usernameText
			up.PostUserId = userId

			up.FriendUserId = dataUserProfile.UserId
			up.FriendUserName = dataUserProfile.UserName

			up.Message = messageBody
			up.RawMessage = messageBodySource
			up.PostTimeVal = messageTimeString
			up.PostTime = messageTime
			up.PostType = data.PostTypeEnum.VisitorPost

			dataVisitorPosts = append(dataVisitorPosts, up)
		}
	})

	return dataUserProfile, dataVisitorPosts, dataFriendList, nil
}

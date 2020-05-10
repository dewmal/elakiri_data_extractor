package data

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type Alias = int
type postType struct {
	VisitorPost Alias
	ThreadPost  Alias
}

var PostTypeEnum = &postType{
	VisitorPost: 1,
	ThreadPost:  2,
}

type VisitorMessage struct {
	gorm.Model
	PostUserId     int64
	FriendUserId   int64
	PostUserName   string
	FriendUserName string
	Message        string
	RawMessage     string
	PostTimeVal    string
	PostTime       time.Time
	PostType       Alias
	ExtractedUrl   string
}

type UserPost struct {
	gorm.Model
	PostId        int64 `gorm:"unique;not null"`
	PostCount     int64
	Username      string
	UserId        int64
	Message       string
	MessageSource string
	PostTimeVal   string
	PostTime      time.Time
	PostType      Alias
	ThreadId      int64
	RelatedPosts  pq.Int64Array `gorm:"type:BigInt[];column:related_posts"`
	ExtractedUrl  string
}

type UserFriend struct {
	UserName     string
	UserId       string
	ExtractedUrl string
}

type UserProfile struct {
	gorm.Model
	UserName       string
	UserId         int64 `gorm:"unique;not null"`
	JoinDateVal    string
	JoinDate       time.Time
	TotalPost      int64
	MemberStatus   string
	ReputationRank int
	TotalPageVisit int64
	Friends        pq.Int64Array `gorm:"type:BigInt[];column:friend_list"`
	LastVisitors   pq.Int64Array `gorm:"type:BigInt[];column:last_visitor_list"`
	Location       string
	ExtractedUrl   string
}

type Thread struct {
	gorm.Model
	OwnerUser    string
	OwnerId      int64
	ThreadId     int64 `gorm:"unique;not null"`
	Title        string
	PostDateTime time.Time
	Categories   pq.StringArray `gorm:"type:VARCHAR(50)[];column:categories"`
	Tags         pq.StringArray `gorm:"type:VARCHAR(50)[];column:tags"`
	Rating       int64
	ExtractedUrl string
}

type ErrorVisitedUrl struct {
	gorm.Model
	VisitedUrl string
}

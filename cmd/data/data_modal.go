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

type UserPost struct {
	gorm.Model
	PostId        string `gorm:"unique;not null"`
	Username      string
	UserId        int64
	Message       string
	MessageSource string
	PostTimeVal   string
	PostTime      time.Time
	PostType      Alias
	ThreadId      int64
}

type UserFriend struct {
	UserName string
	UserId   string
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
}

type Thread struct {
	gorm.Model
	OwnerUser    string
	OwnerId      int64
	ThreadId     int64 `gorm:"unique;not null"`
	Title        string
	PostDateTime time.Time
	Categories   pq.StringArray `gorm:"type:VARCHAR(50)[];column:tags"`
	Tags         pq.StringArray `gorm:"type:VARCHAR(50)[];column:tags"`
	Rating       int64
}

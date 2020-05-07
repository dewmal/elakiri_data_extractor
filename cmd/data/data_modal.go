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
	UserId        string
	Message       string
	MessageSource string
	PostTimeVal   string
	PostTime      time.Time
	PostType      Alias
}

type UserFriend struct {
	UserName string
	UserId   string
}

type UserProfile struct {
	gorm.Model
	UserName       string
	UserId         string `gorm:"unique;not null"`
	JoinDateVal    string
	JoinDate       time.Time
	TotalPost      int64
	MemberStatus   string
	ReputationRank int
	TotalPageVisit int
	Friends        pq.StringArray `gorm:"type:VARCHAR(50)[];column:friend_list"`
	LastVisitors   pq.StringArray `gorm:"type:VARCHAR(50)[];column:last_visitor_list"`
}

type Thread struct {
	gorm.Model
	OwnerUser string
	OwnerId   string
}

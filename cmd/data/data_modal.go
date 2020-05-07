package data

import (
	"github.com/jinzhu/gorm"
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
	Username      string
	UserId        string
	Message       string
	MessageSource string
	PostTimeVal   string
	PostTime      time.Time
	PostType      Alias
}

type UserProfile struct {
	gorm.Model
	UserName    string
	UserId      string
	JoinDateVal string
	JoinDate    time.Time
	TotalPost   int64
}
type UserProfileHasFriend struct {
	gorm.Model
	UserName   string
	UserId     string
	FriendName string
	FriendId   string
}

type Thread struct {
	gorm.Model
	OwnerUser string
	OwnerId   string
}

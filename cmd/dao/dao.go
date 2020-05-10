package dao

import (
	"github.com/jinzhu/gorm"
	"time"
	"webcrawler/cmd/data"
)

func SaveVisitorPost(db *gorm.DB, visitorPost data.VisitorMessage) {
	db.Create(&visitorPost)
}
func SaveThread(db *gorm.DB, thread data.Thread) {
	var dbThread data.Thread
	db.Where(&data.Thread{
		ThreadId: thread.ThreadId,
	}).FirstOrCreate(&dbThread)
	//OwnerUser
	if thread.OwnerUser != "" {
		dbThread.OwnerUser = thread.OwnerUser
	}
	//OwnerId
	if thread.OwnerId != 0 {
		dbThread.OwnerId = thread.OwnerId
	}
	//ThreadId
	if thread.ThreadId != 0 {
		dbThread.ThreadId = thread.ThreadId
	}
	//Title
	if thread.Title != "" {
		dbThread.Title = thread.Title
	}
	//PostDateTime
	if thread.PostDateTime != (time.Time{}) {
		dbThread.PostDateTime = thread.PostDateTime
	}
	//Categories
	if thread.Categories != nil {
		dbThread.Categories = thread.Categories
	}
	//Tags
	if thread.Tags != nil {
		dbThread.Tags = thread.Tags
	}
	//Rating
	if thread.Rating != 0 {
		dbThread.Rating = thread.Rating
	}
	//Extracted URL
	if thread.ExtractedUrl != "" {
		dbThread.ExtractedUrl = thread.ExtractedUrl
	}

	db.Save(&dbThread)
}

func SaveUserPost(db *gorm.DB, userPost data.UserPost) {
	var dbUserPost data.UserPost
	db.Where(&data.UserPost{
		PostId: userPost.PostId,
	}).FirstOrCreate(&dbUserPost)

	if userPost.UserId != 0 {
		dbUserPost.UserId = userPost.UserId
	}
	if userPost.ThreadId != 0 {
		dbUserPost.ThreadId = userPost.ThreadId
	}
	if userPost.PostType != 0 {
		dbUserPost.PostType = userPost.PostType
	}
	if userPost.PostTimeVal != "" {
		dbUserPost.PostTimeVal = userPost.PostTimeVal
	}
	if userPost.MessageSource != "" {
		dbUserPost.MessageSource = userPost.MessageSource
	}
	if userPost.Username != "" {
		dbUserPost.Username = userPost.Username
	}
	if userPost.Message != "" {
		dbUserPost.Message = userPost.Message
	}
	if userPost.PostCount != 0 {
		dbUserPost.PostCount = userPost.PostCount
	}
	if userPost.PostTime != (time.Time{}) {
		dbUserPost.PostTime = userPost.PostTime
	}
	if userPost.RelatedPosts != nil {
		dbUserPost.RelatedPosts = userPost.RelatedPosts
	}
	//Extracted URL
	if userPost.ExtractedUrl != "" {
		dbUserPost.ExtractedUrl = userPost.ExtractedUrl
	}

	db.Save(&dbUserPost)

	//PostCount
	//Username
	//UserId
	//Message
	//MessageSource
	//PostTimeVal
	//PostTime
	//PostType
	//ThreadId
	//RelatedPosts

}

func SaveUserProfile(db *gorm.DB, userProfile data.UserProfile) {
	var dbUserProfile data.UserProfile
	db.Where(&data.UserProfile{
		UserId: userProfile.UserId,
	}).FirstOrCreate(&dbUserProfile)

	if userProfile.UserName != "" {
		dbUserProfile.UserName = userProfile.UserName
	}
	if userProfile.ReputationRank != 0 {
		dbUserProfile.ReputationRank = userProfile.ReputationRank
	}
	if userProfile.Location != "" {
		dbUserProfile.Location = userProfile.Location
	}
	if userProfile.TotalPageVisit != 0 {
		dbUserProfile.TotalPageVisit = userProfile.TotalPageVisit
	}
	if userProfile.LastVisitors != nil {
		dbUserProfile.LastVisitors = userProfile.LastVisitors
	}
	if userProfile.Friends != nil {
		dbUserProfile.Friends = userProfile.Friends
	}
	if userProfile.JoinDate != (time.Time{}) {
		dbUserProfile.JoinDate = userProfile.JoinDate
	}
	if userProfile.JoinDateVal != "" {
		dbUserProfile.JoinDateVal = userProfile.JoinDateVal
	}
	if userProfile.TotalPost != 0 {
		dbUserProfile.TotalPost = userProfile.TotalPost
	}
	if userProfile.TotalPageVisit != 0 {
		dbUserProfile.TotalPageVisit = userProfile.TotalPageVisit
	}
	//Extracted URL
	if userProfile.ExtractedUrl != "" {
		dbUserProfile.ExtractedUrl = userProfile.ExtractedUrl
	}

	db.Save(&dbUserProfile)
}

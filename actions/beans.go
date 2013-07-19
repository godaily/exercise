package actions

import (
	. "code.google.com/p/go.crypto/scrypt"
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

const (
	EXERCISE_MODULE = iota + 1
	QUESTION_MODULE
	ABOUT_MODULE
)

type User struct {
	Id           int64
	LoginName    string `xorm:"unique not null"`
	UserName     string `xorm:"not null"`
	Email        string `xorm:"unique not null"`
	Password     string `xorm:"varchar(128)"`
	NumFollowers int
	NumAsks      int
	NumAnswers   int
	NumComments  int
	NumExercises int
	NumQuestions int
	Avatar       string `xorm:"varchar(2048)"`
}

const (
	USER_ID_TAG     = "UserId"
	USER_NAME_TAG   = "UserName"
	USER_AVATAR_TAG = "UserAvatar"
)

func (u *User) EncodePasswd() error {
	newPasswd, err := Key([]byte(u.Password), []byte("!#@FDEWREWR&*("), 16384, 8, 1, 64)
	u.Password = fmt.Sprintf("%x", newPasswd)
	return err
}

func (u *User) BuildAvatar() {
	m := md5.New()
	m.Write([]byte(strings.ToLower(strings.Trim(u.Email, " "))))
	dg := fmt.Sprintf("%x", m.Sum(nil))
	fmt.Println(dg)
	u.Avatar = fmt.Sprintf("http://www.gravatar.com/avatar/%v", dg)
}

type UserFollow struct {
	UserId     int64
	FollowerId int64
	FollowTime time.Time
}

type Question struct {
	Id           int64
	Type         int    `xorm:"not null"`
	Title        string `xorm:"varchar(500) not null"`
	Creator      User   `xorm:"creator_id int(11)"`
	Content      string `xorm:"text not null"`
	NumComments  int
	NumFollowers int
	NumAnswers   int
	NumReads     int
	NumUps       int
	Created      time.Time `xorm:"not null"`
	LastUpdated  time.Time `xorm:"not null"`
}

type QuestionFollow struct {
	QuestionId int64
	FollowerId int64
	FollowTime time.Time
}

type QuestionComment struct {
	Id          int64
	QuestionId  int64
	Creator     User      `xorm:"creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null"`
	LastUpdated time.Time `xorm:"not null"`
}

type Answer struct {
	Id          int64
	QuestionId  int64
	Creator     User      `xorm:"creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null"`
	LastUpdated time.Time `xorm:"not null"`
	NumComments int
	NumUps      int
}

type AnswerUp struct {
	AnswerId   int64
	UserId     int64
	FollowTime time.Time
}

type AnswerComment struct {
	Id          int64
	AnswerId    int64
	Creator     User      `xorm:"creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null"`
	LastUpdated time.Time `xorm:"not null"`
}

type Tag struct {
	Id    int64
	Name  string `xorm:"unique not null"`
	Total int
}

type QuestionTag struct {
	QuestionId int64
	TagId      int64
}

type Message struct {
	Id       int64
	Sender   User   `xorm:"sender_id int(11)"`
	Receiver User   `xorm:"receiver_id int(11)"`
	Content  string `xorm:"text not null"`
	SendTime time.Time
	ReadTime time.Time
}

type Topic struct {
	Id           int64
	ParentId     int64
	Name         string `xorm:"not null"`
	Url          string `xorm:"varchar(2048) not null"`
	NumFollowers int
	NumQuestions int
	Admin        User `xorm:"admin_id int(11)"`
}

type QuestionTopic struct {
	QuestionId int64
	TopicId    int64
}

type TopicFollow struct {
	TopicId    int64
	FollowerId int64
}

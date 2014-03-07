package actions

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
	"time"

	. "code.google.com/p/go.crypto/scrypt"
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
	Password     string `xorm:"varchar(128) not null"`
	NumFollowers int
	NumAsks      int
	NumAnswers   int
	NumComments  int
	NumExercises int
	NumQuestions int
	Avatar       string `xorm:"varchar(2048) not null"`
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
	u.Avatar = fmt.Sprintf("http://www.gravatar.com/avatar/%v", dg)
}

type UserFollow struct {
	UserId     int64     `xorm:"unique(uf)"`
	FollowerId int64     `xorm:"unique(uf)"`
	FollowTime time.Time `xorm:"not null"`
}

type Question struct {
	Id           int64
	Type         int    `xorm:"not null"`
	Title        string `xorm:"varchar(500) not null"`
	Creator      User   `xorm:"index creator_id int(11)"`
	Content      string `xorm:"text not null"`
	NumComments  int
	NumFollowers int
	NumAnswers   int
	NumReads     int
	NumUps       int
	Created      time.Time `xorm:"not null created"`
	LastUpdated  time.Time `xorm:"not null updated"`
}

type QuestionFollow struct {
	QuestionId int64     `xorm:"unique(qf)"`
	FollowerId int64     `xorm:"unique(qf)"`
	FollowTime time.Time `xorm:"not null"`
}

type QuestionComment struct {
	Id          int64
	QuestionId  int64     `xorm:"index"`
	Creator     User      `xorm:"index creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null created"`
	LastUpdated time.Time `xorm:"not null updated"`
}

type Answer struct {
	Id          int64
	QuestionId  int64     `xorm:"index"`
	Creator     User      `xorm:"index creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null created"`
	LastUpdated time.Time `xorm:"not null updated"`
	NumComments int
	NumUps      int
}

type AnswerUp struct {
	AnswerId int64     `xorm:"unique(au)"`
	UserId   int64     `xorm:"unique(au)"`
	UpTime   time.Time `xorm:"not null"`
}

type AnswerComment struct {
	Id          int64
	AnswerId    int64     `xorm:"index"`
	Creator     User      `xorm:"index creator_id int(11)"`
	Content     string    `xorm:"text not null"`
	Created     time.Time `xorm:"not null created"`
	LastUpdated time.Time `xorm:"not null updated"`
}

type Tag struct {
	Id    int64
	Name  string `xorm:"unique not null"`
	Total int
}

type QuestionTag struct {
	QuestionId int64 `xorm:"unique(qt)"`
	TagId      int64 `xorm:"unique(qt)"`
}

type Message struct {
	Id       int64
	Sender   User   `xorm:"index sender_id int(11)"`
	Receiver User   `xorm:"index receiver_id int(11)"`
	Content  string `xorm:"text not null"`
	SendTime time.Time
	ReadTime time.Time
}

type Topic struct {
	Id           int64
	ParentId     int64  `xorm:"index"`
	Name         string `xorm:"unique not null"`
	Url          string `xorm:"varchar(2048) not null"`
	NumFollowers int
	NumQuestions int
	Admin        User `xorm:"index admin_id int(11)"`
}

type QuestionTopic struct {
	QuestionId int64 `xorm:"unique(qtopic)"`
	TopicId    int64 `xorm:"unique(qtopic)"`
}

type TopicFollow struct {
	TopicId    int64 `xorm:"unique(tf)"`
	FollowerId int64 `xorm:"unique(tf)"`
}

type News struct {
	Id         int64
	Title      string    `xorm:"notnull"`
	Link       string    `xorm:"varchar(2048) notnull"`
	Creator    User      `xorm:"index creator_id int"`
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
	NumComment int
}

func (this *News) Domain() string {
	u, err := url.Parse(this.Link)
	if err != nil {
		return ""
	}
	return u.Host
}

type NewsComment struct {
	Id      int64
	NewsId  int64     `xorm:"index"`
	Content string    `xorm:"text"`
	Creator User      `xorm:"index creator_id"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

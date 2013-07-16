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

type Exercise struct {
	Id         int64
	CreatorId  int64
	Created    time.Time
	Title      string `xorm:"varchar(500)"`
	Content    string `xorm:"text"`
	ShowDate   time.Time
	NumAnswers int
}

type ExerciseAnswer struct {
	Id         int64
	ExerciseId int64
	Created    time.Time
	Content    string `xorm:"text"`
}

type User struct {
	Id           int64
	LoginName    string `xorm:"unique"`
	UserName     string
	Email        string
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

type Question struct {
	Id           int64
	Title        string `xorm:"varchar(200)"`
	Content      string `xorm:"text"`
	NumComments  int
	NumFollowers int
	NumAnswers   int
	NumReads     int
	NumUps       int
	LastUpdated  time.Time
}

type QuestionFollow struct {
	QuestionId int64
	FollowerId int64
}

type UserFollow struct {
	UserId     int64
	FollowerId int64
}

type Answer struct {
	Id          int64
	QuestionId  int64
	Content     string `xorm:"text"`
	LastUpdated time.Time
	NumComments int
	NumAgrees   int
}

type AnswerAgree struct {
	AnswerId int64
	UserId   int64
}

type QuestionComment struct {
	Id         int64
	QuestionId int64
	Content    string `xorm:"text"`
}

type AnswerComment struct {
	Id       int64
	AnswerId int64
	Conetn   string `xorm:"text"`
}

type Tag struct {
	Id    int64
	Name  string `xorm:"unique"`
	Total int
}

type QuestionTag struct {
	QuestionId int64
	TagId      int64
}

type Message struct {
	Id         int64
	SenderId   int64
	ReceiverId int64
	Content    string `xorm:"text"`
}

type Topic struct {
	Id           int64
	ParentId     int64
	Name         string
	Url          string `xorm:"varchar(2048)"`
	NumFollowers int
	NumQuestions int
	AdminId      int64
}

type QuestionTopic struct {
	QuestionId int64
	TopicId    int64
}

type TopicFollow struct {
	TopicId    int64
	FollowerId int64
}

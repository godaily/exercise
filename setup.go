package main

import (
	. "./actions"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	"github.com/lunny/xweb"
	"io/ioutil"
)

var (
	engine *xorm.Engine
)

func main() {
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println("Fail to load configuration:", err)
		return
	}

	cfgs := xweb.SimpleParse(string(data))
	engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println("Fail to connect to database:", err)
		return
	}
	engine.ShowSQL = true

	err = engine.Sync(&User{}, &Question{},
		&QuestionFollow{}, &UserFollow{}, &Answer{}, &AnswerUp{},
		&QuestionComment{}, &AnswerComment{}, &Tag{}, &QuestionTag{},
		&Message{}, &Topic{}, &QuestionTopic{}, &TopicFollow{})

	if err != nil {
		fmt.Println("Fail to create tables:", err)
		return
	}

	fmt.Println("Set up successfully!")
}

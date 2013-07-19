package main

import (
	. "./actions"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/lunny/xorm"
	. "github.com/lunny/xweb"
	"io/ioutil"
)

var (
	engine *Engine
)

func main() {
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := SimpleParse(string(data))
	engine, err = NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	engine.ShowSQL = true

	err = engine.CreateTables(&User{}, &Question{},
		&QuestionFollow{}, &UserFollow{}, &Answer{}, &AnswerUp{},
		&QuestionComment{}, &AnswerComment{}, &Tag{}, &QuestionTag{},
		&Message{}, &Topic{}, &QuestionTopic{}, &TopicFollow{})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Set up successfully!")
}

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
		fmt.Println("Fail to load configuration:", err)
		return
	}

	cfgs := SimpleParse(string(data))
	engine, err = NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println("Fail to connect to database:", err)
		return
	}
	engine.ShowSQL = true

	err = engine.CreateTables(&Exercise{}, &ExerciseAnswer{}, &User{}, &Question{},
		&QuestionFollow{}, &UserFollow{}, &Answer{}, &AnswerAgree{},
		&QuestionComment{}, &AnswerComment{}, &Tag{}, &QuestionTag{},
		&Message{}, &Topic{}, &QuestionTopic{}, &TopicFollow{})

	if err != nil {
		fmt.Println("Fail to create tables:", err)
		return
	}

	fmt.Println("Set up successfully!")
}

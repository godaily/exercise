package main

import (
	. "./actions"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/lunny/xorm"
	. "github.com/lunny/xweb"
)

func main() {
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := SimpleParse(string(data))

	Orm, err = NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	Orm.ShowSQL = true

	AddAction(&MainAction{})
	AddRouter("/exercise/", &ExerciseAction{})
	AddRouter("/question/", &QuestionAction{})
	app := MainServer().RootApp
	loginFilter := NewLoginFilter(app, USER_ID_TAG, "/login")
	loginFilter.AddAskLoginUrls("/exercise/add", "/exercise/sub")
	app.AddFilter(loginFilter)
	Run("0.0.0.0:" + cfgs["port"])
}

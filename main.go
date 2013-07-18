package main

import (
	"./actions"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/lunny/xorm"
	. "github.com/lunny/xweb"
)

const APP_VER = "0.0.1 Beta"

func init() {
	// Setting application version.
	actions.AppVer = "v" + APP_VER
}

func main() {
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := SimpleParse(string(data))

	actions.Orm, err = NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	actions.Orm.ShowSQL = true

	AddAction(&actions.HomeAction{})
	AddRouter("/exercise", &actions.ExerciseAction{})
	AddRouter("/question", &actions.QuestionAction{})
	app := MainServer().RootApp
	loginFilter := NewLoginFilter(app, actions.USER_ID_TAG, "/login")
	loginFilter.AddAskLoginUrls("/exercise/add", "/exercise/sub")
	app.AddFilter(loginFilter)
	Run("0.0.0.0:" + cfgs["port"])
}

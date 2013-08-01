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
	//actions.AppVer = "v" + APP_VER
}

func main() {
	// load config
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := SimpleParse(string(data))

	// create Orm
	actions.Orm, err = NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	actions.Orm.ShowSQL = true

	// add actions
	AddAction(&actions.HomeAction{})
	AddRouter("/exercise", &actions.ExerciseAction{})
	AddRouter("/question", &actions.QuestionAction{})
	AddAction(&actions.UserAction{})

	// add login filter
	app := MainServer().RootApp
	loginFilter := NewLoginFilter(app, actions.USER_ID_TAG, "/login")
	loginFilter.AddAnonymousUrls("/", "/exercise/", "/exercise/compile",
		"/login", "/about")
	app.AddFilter(loginFilter)

	// add func app scope
	app.AddFunc("AppVer", func() string {
		return "v" + APP_VER
	})

	app.SetConfig("Orm", actions.Orm)

	// run the web server
	Run(fmt.Sprintf("%v:%v", cfgs["address"], cfgs["port"]))
}

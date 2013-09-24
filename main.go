package main

import (
	"./actions"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	"github.com/lunny/xweb"
)

const APP_VER = "0.0.1 Beta"

func main() {
	// load config
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := xweb.SimpleParse(string(data))

	// create Orm
	actions.Orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	actions.Orm.ShowSQL = true
	actions.Orm.ShowDebug = true

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	actions.Orm.SetDefaultCacher(cacher)

	// add actions
	xweb.AddAction(&actions.HomeAction{})
	xweb.AddRouter("/exercise", &actions.ExerciseAction{})
	xweb.AddRouter("/question", &actions.QuestionAction{})
	xweb.AddAction(&actions.UserAction{})

	// add login filter
	app := xweb.MainServer().RootApp
	loginFilter := xweb.NewLoginFilter(app, actions.USER_ID_TAG, "/login")
	loginFilter.AddAnonymousUrls("/", "/exercise/", "/exercise/compile",
		"/login", "/about", "/register")
	app.AddFilter(loginFilter)

	// add func app scope
	app.AddFunc("AppVer", func() string {
		return "v" + APP_VER
	})

	app.SetConfig("Orm", actions.Orm)

	// run the web server
	xweb.Run(fmt.Sprintf("%v:%v", cfgs["address"], cfgs["port"]))
}

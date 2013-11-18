package main

import (
	"./actions"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
	"github.com/lunny/xweb"
	//"runtime"
)

const APP_VER = "0.0.2 Beta"

func main() {
	//runtime.GOMAXPROCS(2)

	// load config
	var err error
	data, err := ioutil.ReadFile("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := xweb.SimpleParse(string(data))

	// create Orm
	var orm *xorm.Engine
	orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	orm.ShowSQL = (cfgs["showSQL"] == "1")
	orm.ShowDebug = (cfgs["showDebug"] == "1")

	if cfgs["useCache"] == "1" {
		//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
		orm.SetDefaultCacher(cacher)
	}

	app := xweb.MainServer().RootApp
	app.SetConfig("Orm", orm)

	// add actions
	xweb.AddAction(&actions.HomeAction{})
	xweb.AutoAction(&actions.ExerciseAction{}, &actions.QuestionAction{})
	xweb.AddAction(&actions.UserAction{})

	// add login filter
	loginFilter := xweb.NewLoginFilter(app, actions.USER_ID_TAG, "/login")
	loginFilter.AddAnonymousUrls("/", "/exercise/", "/exercise/compile",
		"/login", "/about", "/register")
	app.AddFilter(loginFilter)

	// add func or var app scope
	app.AddTmplVar("AppVer", func() string {
		return "v" + APP_VER
	})

	// run the web server
	xweb.Run(fmt.Sprintf("%v:%v", cfgs["address"], cfgs["port"]))
}

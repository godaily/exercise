package main

import (
	"fmt"
	//"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/config"
	"github.com/lunny/xorm"
	"github.com/lunny/xweb"

	. "github.com/govc/godaily/actions"
)

const APP_VER = "0.0.2 Beta"

func main() {
	//runtime.GOMAXPROCS(2)

	// load config
	var err error
	cfg, err := config.Load("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgs := cfg.Map()

	// create Orm
	var orm *xorm.Engine
	orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%v:%v@%v/%v?charset=utf8",
		cfgs["dbuser"], cfgs["dbpasswd"], cfgs["dbhost"], cfgs["dbname"]))
	if err != nil {
		fmt.Println(err)
		return
	}
	orm.ShowSQL, _ = cfg.GetBool("showSql")
	orm.ShowDebug, _ = cfg.GetBool("showDebug")

	err = orm.Sync(&User{}, &Question{},
		&QuestionFollow{}, &UserFollow{}, &Answer{}, &AnswerUp{},
		&QuestionComment{}, &AnswerComment{}, &Tag{}, &QuestionTag{},
		&Message{}, &Topic{}, &QuestionTopic{}, &TopicFollow{}, &News{})

	if err != nil {
		fmt.Println(err)
		return
	}

	if useCache, _ := cfg.GetBool("useCache"); useCache {
		fmt.Println("using orm cache system ...")
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
		orm.SetDefaultCacher(cacher)
	}

	app := xweb.RootApp()
	app.SetConfig("Orm", orm)

	// add actions
	xweb.AddAction(&HomeAction{})
	xweb.AutoAction(&ExerciseAction{}, &QuestionAction{}, &NewsAction{})
	xweb.AddAction(&UserAction{})

	// add login filter
	loginFilter := xweb.NewLoginFilter(app, USER_ID_TAG, "/login")
	loginFilter.AddAnonymousUrls("/", "/exercise", "/exercise/compile",
		"/news", "/login", "/about", "/register")
	app.AddFilter(loginFilter)

	// add func or var app scope
	app.AddTmplVar("AppVer", func() string {
		return "v" + APP_VER
	})

	// run the web server
	xweb.Run(fmt.Sprintf("%v:%v", cfgs["address"], cfgs["port"]))
}

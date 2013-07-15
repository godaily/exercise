package actions

import (
	//"fmt"
	. "github.com/lunny/xorm"
	. "github.com/lunny/xweb"
	//. "xweb"
)

var (
	Orm *Engine
)

type MainAction struct {
	BaseAction

	root     Mapper `xweb:"/"`
	about    Mapper
	register Mapper
	login    Mapper

	User       User
	Message    string
	RePassword string
}

func (c *MainAction) Init() {
	c.BaseAction.Init()
	c.AddFunc("isCurModule", c.IsCurModule)
}

func (c *MainAction) IsCurModule(cur int) bool {
	return ABOUT_MODULE == cur
}

func (c *MainAction) About() {
	c.Render("about.html")
}

func (c *MainAction) Root() {
	c.Go("root", &ExerciseAction{})
}

func (c *MainAction) Login() error {
	if c.Method() == "GET" {
		return c.Render("login.html")
	} else if c.Method() == "POST" {
		c.User.EncodePasswd()
		has, err := Orm.Get(&c.User)
		if err == nil {
			if has {
				c.SetSession(USER_ID_TAG, c.User.Id)
				return c.Go("root")
			}
			return c.Go("login?message=账号或密码错误")
		}
		return err
	}
	return NotSupported()
}

func (c *MainAction) Register() error {
	if c.Method() == "GET" {
		return c.Render("register.html")
	} else if c.Method() == "POST" {
		if c.RePassword != c.User.Password {
			return c.Go("register?message=两次密码不匹配")
		}
		c.User.EncodePasswd()
		_, err := Orm.Insert(&c.User)
		if err == nil {
			return c.Render("registerok.html")
		}
		return err
	}
	return NotSupported()
}

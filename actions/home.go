package actions

import (
	. "github.com/lunny/xorm"
	"github.com/lunny/xweb"
	"regexp"
)

var (
	Orm    *Engine
	AppVer string
)

type HomeAction struct {
	BaseAction

	root     xweb.Mapper `xweb:"/"`
	about    xweb.Mapper
	register xweb.Mapper
	login    xweb.Mapper
	logout   xweb.Mapper

	User       User
	Message    string
	RePassword string
}

func (c *HomeAction) Init() {
	c.BaseAction.Init()
}

func (c *HomeAction) About() {
	c.Render("about.html", &xweb.T{
		"IsAbout": true,
	})
}

func (c *HomeAction) Root() error {
	return c.Go("root", &ExerciseAction{})
	/*return c.Render("home/root.html", &xweb.T{
		"IsHome": true,
	})*/
}

func (c *HomeAction) Login() error {
	if c.Method() == "GET" {
		return c.Render("login.html")
	} else if c.Method() == "POST" {
		if len(c.User.LoginName) <= 0 {
			return c.Go("login?message=登录名不能为空")
		}
		c.User.EncodePasswd()
		has, err := Orm.Get(&c.User)
		if err == nil {
			if has {
				c.SetSession(USER_ID_TAG, c.User.Id)
				c.SetSession(USER_NAME_TAG, c.User.UserName)
				c.SetSession(USER_AVATAR_TAG, c.User.Avatar)
				return c.Go("root")
			}
			return c.Go("login?message=账号或密码错误")
		}
		return err
	}
	return xweb.NotSupported()
}

func (c *HomeAction) Logout() error {
	c.DelSession(USER_ID_TAG)
	c.DelSession(USER_NAME_TAG)
	c.DelSession(USER_AVATAR_TAG)
	return c.Go("root")
}

func (c *HomeAction) Register() error {
	if c.Method() == "GET" {
		return c.Render("register.html")
	} else if c.Method() == "POST" {
		if len(c.User.LoginName) <= 0 {
			return c.Go("register?message=登录名不能为空")
		}
		if len(c.User.Email) <= 0 {
			return c.Go("register?message=email地址不能为空")
		}
		isEmail, err := regexp.MatchString("\\w*@\\w*.\\w*", c.User.Email)
		if err != nil {
			return err
		}
		if !isEmail {
			return c.Go("register?message=email地址不正确")
		}
		if len(c.User.Password) <= 0 {
			return c.Go("register?message=密码不能为空")
		}
		if c.RePassword != c.User.Password {
			return c.Go("register?message=两次密码不匹配")
		}
		u := &User{}
		has, err := Orm.Sql("select * from user where login_name=? or email =?",
			c.User.LoginName, c.User.Email).Get(u)
		if has {
			return c.Go("register?message=登录名或者email地址重复")
		}
		if err != nil {
			return err
		}

		c.User.EncodePasswd()
		c.User.BuildAvatar()
		_, err = Orm.Insert(&c.User)
		if err == nil {
			return c.Render("registerok.html")
		}
		return err
	}
	return xweb.NotSupported()
}

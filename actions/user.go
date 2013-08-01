package actions

import (
	"errors"
	"github.com/lunny/xweb"
)

type UserAction struct {
	BaseAction

	changePass xweb.Mapper `xweb:"/setttings/pass"`
	root       xweb.Mapper `xweb:"/(.*)"`

	Password   string
	Repassword string
	Message    string
}

func (c *UserAction) Init() {
	c.BaseAction.Init()
}

func (c *UserAction) ChangePass() error {
	if c.Method() == "GET" {
		return c.Render("user/pass.html")
	} else if c.Method() == "POST" {
		if c.Password != c.Repassword {
			c.Message = "两次输入密码不匹配"
			return c.Render("user/pass.html")
		}

		user := &User{Password: c.Password}
		err := user.EncodePasswd()
		if err == nil {
			_, err = c.Orm().Id(c.GetLoginUserId()).Update(user)
			//_, err = Orm.Id(c.GetLoginUserId()).Update(user)
			if err == nil {
				return c.Render("user/passok.html")
			}
		}
		return err
	}
	return xweb.NotSupported()
}

func (c *UserAction) Root(name string) error {
	user := &User{LoginName: name}
	has, err := Orm.Get(user)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("不存在的用户")
	}

	return c.Render("user.html", &xweb.T{
		"user": user,
	})
}

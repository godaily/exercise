package actions

import (
	"errors"
	"github.com/lunny/xweb"
)

type UserAction struct {
	BaseAction

	root xweb.Mapper `xweb:"/(.*)"`
}

func (c *UserAction) Init() {
	c.BaseAction.Init()
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

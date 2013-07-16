package actions

import (
	. "github.com/lunny/xweb"
)

type BaseAction struct {
	Action
}

func (c *BaseAction) Init() {
	c.AddFunc("IsLogedIn", c.IsLogedIn)
}

func (c *BaseAction) IsLogedIn() bool {
	return c.GetSession(USER_ID_TAG) != nil
}

func (c *BaseAction) GetLoginUserId() int64 {
	return c.GetSession(USER_ID_TAG).(int64)
}

func (c *BaseAction) GetLoginUser() *User {
	return nil
}

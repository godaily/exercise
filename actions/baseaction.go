package actions

import (
	. "github.com/lunny/xweb"
)

type BaseAction struct {
	Action
}

func (c *BaseAction) Init() {
	c.AddFunc("IsLogedIn", c.IsLogedIn)
	c.AddFunc("GetLoginUserAvatar", c.GetLoginUserAvatar)
	c.AddFunc("GetLoginUserName", c.GetLoginUserName)
}

func (c *BaseAction) IsLogedIn() bool {
	return c.GetSession(USER_ID_TAG) != nil
}

func (c *BaseAction) GetLoginUserId() int64 {
	return c.GetSession(USER_ID_TAG).(int64)
}

func (c *BaseAction) GetLoginUserName() string {
	return c.GetSession(USER_NAME_TAG).(string)
}

func (c *BaseAction) GetLoginUserAvatar() string {
	return c.GetSession(USER_AVATAR_TAG).(string)
}

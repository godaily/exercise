package actions

import (
	"github.com/lunny/xorm"
	"github.com/lunny/xweb"
)

type BaseAction struct {
	xweb.Action
	Orm *xorm.Engine
}

func (c *BaseAction) Init() {
	c.AddTmplVars(&xweb.T{"IsLogedIn": c.IsLogedIn,
		"GetLoginUserAvatar": c.GetLoginUserAvatar,
		"GetLoginUserName":   c.GetLoginUserName,
		"GetLoginUserId":     c.GetLoginUserId,
	})
	c.Orm = c.App.GetConfig("Orm").(*xorm.Engine)
}

func (c *BaseAction) IsLogedIn() bool {
	return c.GetSession(USER_ID_TAG) != nil
}

func (c *BaseAction) GetLoginUserId() int64 {
	id := c.GetSession(USER_ID_TAG)
	if id != nil {
		return id.(int64)
	}
	return 0
}

func (c *BaseAction) GetLoginUserName() string {
	name := c.GetSession(USER_NAME_TAG)
	if name != nil {
		return name.(string)
	}
	return ""
}

func (c *BaseAction) GetLoginUserAvatar() string {
	avatar := c.GetSession(USER_AVATAR_TAG)
	if avatar != nil {
		return avatar.(string)
	}
	return ""
}

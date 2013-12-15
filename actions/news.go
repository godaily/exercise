package actions

import (
    "github.com/lunny/xweb"
)

type NewsAction struct {
    BaseAction

    root    xweb.Mapper `xweb:"/"`
    add     xweb.Mapper
    edit    xweb.Mapper
    comment xweb.Mapper
    clicked xweb.Mapper

    Id int64
}

func (this *NewsAction) Root() error {
    news := make([]News, 0)
    err := this.Orm.Desc("(id)").Find(news)
    if err == nil {
        err = this.Render("news/root.html", &xweb.T{
            "news": news,
        })
    }
    return err
}

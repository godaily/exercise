package actions

import (
    "time"

    "github.com/lunny/xweb"    
)

type QuestionAction struct {
    BaseAction

    root     xweb.Mapper `xweb:"/"`
    ask      xweb.Mapper
    answer   xweb.Mapper
    question xweb.Mapper `xweb:"/q"`

    TheQuestion Question
    TheAnswer   Answer
    QuestionId  int64
}

func (c *QuestionAction) Init() {
    c.BaseAction.Init()
}

func (c *QuestionAction) Root() error {
    questions := make([]Question, 0)
    err := c.Orm.Find(&questions)
    if err == nil {
        return c.Render("question/root.html", &xweb.T{
            "questions": &questions,
        })
    }
    return err
}

func (c *QuestionAction) Question() error {
    answers := make([]Answer, 0)
    _, err := c.Orm.Id(c.QuestionId).Get(&c.TheQuestion)
    if err != nil {
        return err
    }

    err = c.Orm.Find(&answers, &Answer{QuestionId: c.QuestionId})
    if err == nil {
        return c.Render("question/question.html", &xweb.T{
            "answers": &answers,
        })
    }
    return err
}

func (c *QuestionAction) Ask() error {
    if c.Method() == "GET" {
        return c.Render("question/ask.html")
    } else if c.Method() == "POST" {
        c.TheQuestion.LastUpdated = time.Now()
        _, err := c.Orm.Insert(&c.TheQuestion)
        if err == nil {
            return c.Render("question/askok.html")
        }
        return err
    }
    return xweb.NotSupported()
}

func (c *QuestionAction) Answer() error {
    if c.Method() == "GET" {
        return c.Render("question/answer.html")
    } else if c.Method() == "POST" {
        _, err := c.Orm.Insert(&c.TheAnswer)
        if err == nil {
            return c.Render("question/answerok.html")
        }
        return err
    }
    return xweb.NotSupported()
}

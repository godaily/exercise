package actions

import (
	"time"

	. "github.com/lunny/play-sdk"
	. "github.com/lunny/xweb"
)

type ExerciseAction struct {
	BaseAction

	root        Mapper `xweb:"/"`
	add         Mapper
	sub         Mapper
	compile     Mapper
	addQComment Mapper `xweb:"POST"`
	addAComment Mapper `xweb:"POST"`
	upAnswer    Mapper

	Exercise Question
	Answer   Answer
	QComment QuestionComment
	AComment AnswerComment
	Id       int64
}

var badges []string = []string{
	"important",
	"success",
	"warning",
	"info",
	"inverse",
}

func GetBadge(i int) string {
	return badges[i]
}

func (c *ExerciseAction) Init() {
	c.BaseAction.Init()
	c.AddFunc("getBadge", GetBadge)
	c.AddVar("IsExer", true)
}

func (c *ExerciseAction) UpAnswer() {
	if c.Id > 0 {
		au := &AnswerUp{AnswerId: c.Id, UserId: c.GetLoginUserId()}
		has, err := Orm.Get(au)
		if err == nil {
			if !has {
				_, err = Orm.Insert(au)
				if err == nil {
					c.ServeJson(map[string]interface{}{"res": 1})
					return
				}
			} else {
				c.ServeJson(map[string]interface{}{"res": 2})
				return
			}
		}
		c.ServeJson(map[string]interface{}{"res": 0, "error": err.Error()})
		return
	}
	c.ServeJson(map[string]interface{}{"res": 0, "error": "无效参数"})
	return
}

func (c *ExerciseAction) Add() error {
	if c.Method() == "GET" {
		recentExercises := make([]Question, 0)
		err := Orm.OrderBy("created desc").Limit(5).Find(&recentExercises)
		if err == nil {
			return c.Render("exercise/add.html", &T{
				"exercises": &recentExercises,
			})
		}
		return err
	} else if c.Method() == "POST" {
		c.Exercise.Creator.Id = c.BaseAction.GetLoginUserId()
		c.Exercise.Created = time.Now()
		c.Exercise.Type = EXERCISE_MODULE
		session := Orm.NewSession()
		defer session.Close()
		err := session.Begin()
		if err == nil {
			_, err = session.Insert(&c.Exercise)
			if err == nil {
				_, err = session.Exec("update user set num_questions = num_questions+1 where id = ?", c.GetLoginUserId())
				if err == nil {
					err = session.Commit()
				}
				if err == nil {
					return c.Render("exercise/addok.html")
				}
			}
		}
		if err != nil {
			session.Rollback()
		}
		return err
	}
	return NotSupported()
}

func (c *ExerciseAction) Compile() {
	res, err := Compile(c.Answer.Content)
	if err == nil {
		if res.Errors == "" {
			c.ServeJson(res.Events)
		} else {
			c.ServeJson(map[string]interface{}{"errors": res.Errors})
		}
	} else {
		c.ServeJson(map[string]interface{}{"errors": err.Error()})
	}
}

func (c *ExerciseAction) Sub() error {
	if c.Method() == "GET" {
		return c.Render("exercise/sub.html")
	} else if c.Method() == "POST" {
		if c.Answer.Id == 0 {
			session := Orm.NewSession()
			defer session.Close()
			c.Answer.Created = time.Now()
			c.Answer.Creator.Id = c.GetLoginUserId()
			err := session.Begin()
			if err == nil {
				_, err = session.Insert(&c.Answer)
				if err == nil {
					_, err = session.Exec("update user set num_exercises=num_exercises+1 where id=?", c.GetLoginUserId())
				}
				if err == nil {
					_, err = session.Exec("update question set num_answers=num_answers+1 where id=?", c.Id)
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				session.Rollback()
			}

			if err == nil {
				return c.Render("exercise/subok.html")
			}
			return err
		} else {
			c.Answer.Created = time.Now()
			_, err := Orm.Id(c.Answer.Id).Update(&c.Answer)
			if err == nil {
				return c.Render("exercise/subok.html")
			}
			return err
		}
	}
	return NotSupported()
}

func (c *ExerciseAction) Root() error {
	var has bool
	var err error
	if c.Id == 0 {
		has, err = Orm.OrderBy("created desc").Get(&c.Exercise)
	} else {
		has, err = Orm.Id(c.Id).Get(&c.Exercise)
	}
	if err == nil {
		var answers []Answer
		var qusers, eusers []User
		var curAnswer Answer
		var hasSubmited bool
		var pre, last Question
		var preId, lastId int
		var qcomments []QuestionComment
		var acomments map[int64][]AnswerComment = make(map[int64][]AnswerComment)
		if has {
			_, err = Orm.OrderBy("id desc").Where("id < ?", c.Exercise.Id).Get(&pre)
			if err != nil {
				return err
			}
			preId = int(pre.Id)
			_, err = Orm.OrderBy("id asc").Where("id > ?", c.Exercise.Id).Get(&last)
			if err != nil {
				return err
			}

			err = Orm.OrderBy("created asc").Find(&qcomments, &QuestionComment{QuestionId: c.Exercise.Id})
			if err != nil {
				return err
			}

			lastId = int(last.Id)
			err = Orm.OrderBy("num_ups desc").Find(&answers, &Answer{QuestionId: c.Exercise.Id})
			if err != nil {
				return err
			}

			for _, answer := range answers {
				var ac []AnswerComment
				err = Orm.OrderBy("created asc").Find(&ac, &AnswerComment{AnswerId: answer.Id})
				if err != nil {
					return err
				}
				acomments[answer.Id] = ac
			}

			err = Orm.OrderBy("num_questions desc").Limit(5).Find(&qusers)
			if err != nil {
				return err
			}
			err = Orm.OrderBy("num_exercises desc").Limit(5).Find(&eusers)
			if err != nil {
				return err
			}
			if c.IsLogedIn() {
				curAnswer.Creator.Id = c.GetLoginUserId()
				curAnswer.QuestionId = c.Exercise.Id
				hasSubmited, err = Orm.Get(&curAnswer)
				if err != nil {
					return err
				}
			}
		}
		return c.Render("exercise/root.html", &T{
			"has":         has,
			"preId":       preId,
			"lastId":      lastId,
			"answers":     &answers,
			"qusers":      &qusers,
			"eusers":      &eusers,
			"hasSubmited": hasSubmited,
			"curAnswer":   curAnswer,
			"qcomments":   qcomments,
			"acomments":   acomments,
		})
	}
	return err
}

func (c *ExerciseAction) AddQComment() error {
	if c.Id == 0 {
		c.QComment.Creator.Id = c.GetLoginUserId()
		c.QComment.Created = time.Now()
		c.QComment.LastUpdated = time.Now()
		_, err := Orm.Insert(&c.QComment)
		return err
	} else {
		c.QComment.LastUpdated = time.Now()
		_, err := Orm.Id(c.QComment.Id).Update(&c.QComment)
		return err
	}
}

func (c *ExerciseAction) AddAComment() error {
	if c.Id == 0 {
		c.AComment.Creator.Id = c.GetLoginUserId()
		c.AComment.Created = time.Now()
		c.AComment.LastUpdated = time.Now()
		_, err := Orm.Insert(&c.AComment)
		return err
	} else {
		c.AComment.LastUpdated = time.Now()
		_, err := Orm.Id(c.AComment.Id).Update(&c.AComment)
		return err
	}
}

package actions

import (
	"errors"
	"fmt"
	"time"

	. "github.com/lunny/play-sdk"
	"github.com/lunny/xweb"
)

type ExerciseAction struct {
	BaseAction

	root        xweb.Mapper `xweb:"/"`
	add         xweb.Mapper
	edit        xweb.Mapper
	sub         xweb.Mapper
	addQComment xweb.Mapper `xweb:"POST"`
	delQComment xweb.Mapper
	addAComment xweb.Mapper `xweb:"POST"`
	delAComment xweb.Mapper
	upAnswer    xweb.Mapper
	compile     xweb.Mapper

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
	c.AddTmplVars(&xweb.T{"getBadge": GetBadge,
		"IsExer": true,
	})
}

func (c *ExerciseAction) UpAnswer() {
	if c.Id > 0 {
		au := &AnswerUp{AnswerId: c.Id, UserId: c.GetLoginUserId()}
		has, err := c.Orm.Get(au)
		if err == nil {
			if !has {
				session := c.Orm.NewSession()
				defer session.Close()
				err = session.Begin()
				if err == nil {
					au.UpTime = time.Now()
					_, err = session.Insert(au)
				}
				if err == nil {
					answer := new(Answer)
					has, err = session.Id(c.Id).Get(answer)
					if err == nil {
						if has {
							answer.NumUps += 1
							_, err = session.Cols("num_ups").Id(c.Id).Update(answer)
						} else {
							err = errors.New("answer is not exist.")
						}
					}
				}
				if err != nil {
					session.Rollback()
				} else {
					err = session.Commit()
				}
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
		err := c.Orm.Desc("created").Limit(5).Find(&recentExercises)
		if err == nil {
			return c.Render("exercise/add.html", &xweb.T{
				"exercises": &recentExercises,
			})
		}
		return err
	} else if c.Method() == "POST" {
		c.Exercise.Creator.Id = c.GetLoginUserId()
		c.Exercise.Type = EXERCISE_MODULE
		session := c.Orm.NewSession()
		defer session.Close()
		err := session.Begin()
		if err == nil {
			_, err = session.Insert(&c.Exercise)
			if err == nil {
				_, err = session.Exec("update user set num_questions=num_questions+1 where id = ?", c.GetLoginUserId())
				session.Engine.ClearCacheBean(new(User), c.GetLoginUserId())
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
	return xweb.NotSupported()
}

func (c *ExerciseAction) Edit() error {
	if c.Method() == "GET" {
		if c.Id > 0 {
			has, err := c.Orm.Id(c.Id).Get(&c.Exercise)
			if err != nil {
				return err
			}
			if has {
				if c.Exercise.Creator.Id != c.GetLoginUserId() {
					return c.Go("root")
				}
				recentExercises := make([]Question, 0)
				err := c.Orm.Desc("created").Limit(5).Find(&recentExercises)
				if err == nil {
					return c.Render("exercise/edit.html", &xweb.T{
						"exercises": &recentExercises,
					})
				}
				return err
			}
		}

		return errors.New("参数错误")
	} else if c.Method() == "POST" {
		_, err := c.Orm.Id(c.Exercise.Id).Update(&c.Exercise)
		if err == nil {
			return c.Render("exercise/editok.html")
		}
		return err
	}
	return xweb.NotSupported()
}

func (c *ExerciseAction) Compile() {
	fmt.Println(c.Answer.Content)
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
			session := c.Orm.NewSession()
			defer session.Close()
			c.Answer.Creator.Id = c.GetLoginUserId()
			err := session.Begin()
			if err == nil {
				_, err = session.Insert(&c.Answer)
				if err == nil {
					user := new(User)
					has, err := session.Id(c.GetLoginUserId()).Get(user)
					if err == nil {
						if has {
							_, err = session.Table(user).Id(c.GetLoginUserId()).Update(map[string]interface{}{"num_exercises": user.NumExercises + 1})
						} else {
							err = errors.New("user is not exist.")
						}
					}
				}
				if err == nil {
					question := new(Question)
					has, err := session.Id(c.Id).Get(question)
					if err == nil {
						if has {
							_, err = session.Table(question).Id(c.Id).Update(map[string]interface{}{"num_answers": question.NumAnswers + 1})
						} else {
							err = errors.New("question is not exist.")
						}
					}
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
			_, err := c.Orm.Id(c.Answer.Id).Update(&c.Answer)
			if err == nil {
				return c.Render("exercise/subok.html")
			}
			return err
		}
	}
	return xweb.NotSupported()
}

func (c *ExerciseAction) Root() error {
	var has bool
	var err error
	if c.Id == 0 {
		has, err = c.Orm.Desc("created").Get(&c.Exercise)
	} else {
		has, err = c.Orm.Id(c.Id).Get(&c.Exercise)
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
			_, err = c.Orm.Desc("id").Where("id < ?", c.Exercise.Id).Get(&pre)
			if err != nil {
				return err
			}
			preId = int(pre.Id)
			_, err = c.Orm.Asc("id").Where("id > ?", c.Exercise.Id).Get(&last)
			if err != nil {
				return err
			}

			err = c.Orm.Asc("created").Find(&qcomments, &QuestionComment{QuestionId: c.Exercise.Id})
			if err != nil {
				return err
			}

			lastId = int(last.Id)
			err = c.Orm.Desc("num_ups").Find(&answers, &Answer{QuestionId: c.Exercise.Id})
			if err != nil {
				return err
			}

			for _, answer := range answers {
				var ac []AnswerComment
				err = c.Orm.Asc("created").Find(&ac, &AnswerComment{AnswerId: answer.Id})
				if err != nil {
					return err
				}
				acomments[answer.Id] = ac
			}

			err = c.Orm.Desc("num_questions").Limit(5).Find(&qusers)
			if err != nil {
				return err
			}
			err = c.Orm.Desc("num_exercises").Limit(5).Find(&eusers)
			if err != nil {
				return err
			}
			if c.IsLogedIn() {
				curAnswer.Creator.Id = c.GetLoginUserId()
				curAnswer.QuestionId = c.Exercise.Id
				hasSubmited, err = c.Orm.Get(&curAnswer)
				if err != nil {
					return err
				}
			}
		}
		return c.Render("exercise/root.html", &xweb.T{
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

func (c *ExerciseAction) AddQComment() {
	var err error
	if c.Id == 0 {
		c.QComment.Creator.Id = c.GetLoginUserId()
		_, err = c.Orm.Insert(&c.QComment)
	} else {
		_, err = c.Orm.Id(c.QComment.Id).Update(&c.QComment)
	}
	if err == nil {
		c.ServeJson(map[string]interface{}{"error": ""})
	} else {
		c.ServeJson(map[string]interface{}{"error": err.Error()})
	}
}

func (c *ExerciseAction) DelQComment() {
	if c.Id > 0 {
		q := &QuestionComment{Creator: User{Id: c.GetLoginUserId()}}
		_, err := c.Orm.Id(c.Id).Delete(q)
		if err == nil {
			c.ServeJson(map[string]interface{}{"error": ""})
		} else {
			c.ServeJson(map[string]interface{}{"error": err.Error()})
		}
	} else {
		c.ServeJson(map[string]interface{}{"error": "参数不正确"})
	}
}

func (c *ExerciseAction) AddAComment() {
	var err error
	if c.Id == 0 {
		c.AComment.Creator.Id = c.GetLoginUserId()
		_, err = c.Orm.Insert(&c.AComment)
	} else {
		_, err = c.Orm.Id(c.AComment.Id).Update(&c.AComment)
	}
	if err == nil {
		c.ServeJson(map[string]interface{}{"error": ""})
	} else {
		c.ServeJson(map[string]interface{}{"error": err.Error()})
	}
}

func (c *ExerciseAction) DelAComment() {
	if c.Id > 0 {
		a := &AnswerComment{Creator: User{Id: c.GetLoginUserId()}}
		_, err := c.Orm.Id(c.Id).Delete(a)
		if err == nil {
			c.ServeJson(map[string]interface{}{"error": ""})
		} else {
			c.ServeJson(map[string]interface{}{"error": err.Error()})
		}
	} else {
		c.ServeJson(map[string]interface{}{"error": "参数不正确"})
	}
}

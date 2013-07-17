package actions

import (
	. "github.com/lunny/play-sdk"
	. "github.com/lunny/xweb"
	"time"
)

type ExerciseAction struct {
	BaseAction

	root    Mapper `xweb:"/"`
	add     Mapper
	sub     Mapper
	compile Mapper

	Exercise Exercise
	Answer   ExerciseAnswer
	Id       int64
}

func (c *ExerciseAction) Init() {
	c.BaseAction.Init()
	c.AddFunc("isCurModule", c.IsCurModule)
}

func (c *ExerciseAction) IsCurModule(cur int) bool {
	return EXERCISE_MODULE == cur
}

func (c *ExerciseAction) Add() error {
	if c.Method() == "GET" {
		recentExercises := make([]Exercise, 0)
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
		c.Exercise.ShowDate = time.Now()
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
					_, err = session.Exec("update exercise set num_answers=num_answers+1 where id=?", c.Id)
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				err = session.Rollback()
			}

			if err == nil {
				return c.Render("exercise/subok.html")
			}
			return err
		} else {
			c.Answer.Created = time.Now()
			_, err := Orm.Update(&c.Answer)
			if err == nil {
				return c.Render("exercise/subok.html")
			}
			return err
		}
	}
	return NotSupported()
}

func (c *ExerciseAction) Root() error {
	var preId, lastId int64
	has, err := Orm.Cascade(false).OrderBy("created desc").Get(&c.Exercise)
	if err == nil {
		var answers []ExerciseAnswer
		var qusers, eusers []User
		var curAnswer ExerciseAnswer
		var hasSubmited bool
		if has {
			err = Orm.Find(&answers)
			if err != nil {
				return err
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
				curAnswer.ExerciseId = c.Exercise.Id
				hasSubmited, err = Orm.Get(&curAnswer)
				if err != nil {
					return err
				}
			}
		}
		return c.Render("exercise/root.html", &T{
			"has":         has,
			"answers":     &answers,
			"qusers":      &qusers,
			"eusers":      &eusers,
			"hasSubmited": hasSubmited,
			"curAnswer":   curAnswer,
		})
	}
	return err
}

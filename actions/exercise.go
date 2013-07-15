package actions

import (
	//. "github.com/lunny/xorm"
	. "github.com/lunny/xweb"
	"time"
	//. "xweb"
)

type ExerciseAction struct {
	BaseAction

	root Mapper `xweb:"/"`
	add  Mapper
	sub  Mapper

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
		c.Exercise.CreatorId = c.BaseAction.GetLoginUserId()
		c.Exercise.Created = time.Now()
		_, err := Orm.Insert(&c.Exercise)
		if err == nil {
			return c.Render("exercise/addok.html")
		}
		return err
	}
	return NotSupported()
}

func (c *ExerciseAction) Sub() error {
	if c.Method() == "GET" {
		return c.Render("exercise/sub.html")
	} else if c.Method() == "POST" {
		_, err := Orm.Insert(&c.Answer)
		if err == nil {
			return c.Render("exercise/subok.html")
		}
		return err
	}
	return NotSupported()
}

func (c *ExerciseAction) Root() error {
	has, err := Orm.Get(&c.Exercise)
	if err == nil {
		var answers []ExerciseAnswer
		var qusers []User
		var eusers []User
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
		}
		return c.Render("exercise/root.html", &T{
			"has":     has,
			"answers": &answers,
			"qusers":  &qusers,
			"eusers":  &eusers,
		})
	}
	return err
}

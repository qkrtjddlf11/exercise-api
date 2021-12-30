package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodayExercise struct {
	Id           int     `json:"id"`
	Trainer_Id   int     `json:"trainer_id"`
	Group_Id     int     `json:"group_id"`
	Exercises    string  `json:"exercises"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User int     `json:"created_user"`
	Updated_User int     `json:"updated_user"`
	User_Id      int     `json:"user_id"`
}

func (td TodayExercise) todayExerciseInsertQuery(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_today_exercises(trainer_id, group_id, exercises, created_user, updated_user, user_id) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := stmt.Exec(
		td.Trainer_Id, td.Group_Id, td.Exercises, td.Created_User, td.Updated_User, td.User_Id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	Id = int(id)

	return
}

func TodayExerciseRouter(router *gin.Engine, db *sql.DB) {
	tdExercises := router.Group("/api/t/exercises")

	tdExercises.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Ok!")
	})

	tdExercises.POST("/", func(c *gin.Context) {
		tdExercises := TodayExercise{}
		err := c.Bind(&tdExercises)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		Id, err := tdExercises.todayExerciseInsertQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("ID : %d, Today exercises Successfully Created", Id),
			})
		}
	})
}

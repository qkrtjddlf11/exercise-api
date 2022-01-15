package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type TodayExercise struct {
	Seq           int     `json:"seq"`
	Trainer_Id    string  `json:"trainer_id"`
	Group_Name    string  `json:"group_name"`
	Exercises     string  `json:"exercises"`
	Created_Date  *string `json:"created_date"`
	Updated_Date  *string `json:"updated_date"`
	Created_User  string  `json:"created_user"`
	Updated_User  string  `json:"updated_user"`
	User_Id       string  `json:"user_id"`
	Exercise_Date *string `json:"exercise_date"`
}

func (td TodayExercise) todayExerciseInsertQuery(db *sql.DB) (Seq int, err error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_today_exercises(trainer_id, group_name, exercises, created_user, updated_user, user_id, exercise_date) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := stmt.Exec(
		td.Trainer_Id, td.Group_Name, td.Exercises, td.Created_User, td.Updated_User, td.User_Id, td.Exercise_Date)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	Seq = int(id)

	return
}

func (td TodayExercise) todayExerciseSelectQuery(db *sql.DB) (tdExercises []TodayExercise, err error) {
	rows, err := db.Query(
		`SELECT 
			trainer_id, 
			group_name, 
			exercises, 
			created_date, 
			updated_date,
			created_user,
			updated_user,
			user_id,
			exercise_date FROM t_today_exercises WHERE trainer_id = ? AND group_name = ? `, td.Trainer_Id, td.Group_Name)
	if err != nil {
		return
	}

	for rows.Next() {
		var tdExercise TodayExercise
		rows.Scan(
			&tdExercise.Trainer_Id,
			&tdExercise.Group_Name,
			&tdExercise.Exercises,
			&tdExercise.Created_Date,
			&tdExercise.Updated_Date,
			&tdExercise.Created_User,
			&tdExercise.Updated_User,
			&tdExercise.User_Id,
			&tdExercise.Exercise_Date)
		tdExercises = append(tdExercises, tdExercise)
	}
	defer rows.Close()

	return
}

func postTodayExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		tdExercises := TodayExercise{}
		tdExercises.Trainer_Id, tdExercises.Group_Name = trainer_id, group_name
		if err := c.Bind(&tdExercises); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, tdExercises))
			return
		}

		_, err := tdExercises.todayExerciseInsertQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, tdExercises))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(tdExercises))
		}
	}

	return resultFunc
}

func getTodayExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		tdExercise := TodayExercise{}
		tdExercise.Trainer_Id, tdExercise.Group_Name = trainer_id, group_name
		tdExercises, err := tdExercise.todayExerciseSelectQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, tdExercises))
		}

		c.JSON(http.StatusOK, common.SucceedResponse(tdExercises))
	}

	return resultFunc
}

func TodayExerciseRouter(router *gin.Engine, db *sql.DB) {
	tdExercises := router.Group("/api/t/exercises")

	// curl http://127.0.0.1:8080/api/t/exercises?trainer_id=Park&group_name=dygym -X POST -d '{"exercises": "스쿼트 20회 5Set", "created_user": "Lee", "updated_user": "Lee", "user_id": "Customer2", "exercise_date": "2022-01-18"}' -H "Content-Type: application/json"
	tdExercises.POST("/", postTodayExercise(db))

	// Get All TodayExercises with specific trainer_id and group_name
	// curl http://127.0.0.1:8080/api/t/exercises?trainer_id=Park&group_name=dygym -X GET
	tdExercises.GET("/", getTodayExercise(db))
}

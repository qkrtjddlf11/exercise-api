package routes

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	User_Name     string  `json:"user_name"`
	Exercise_Date *string `json:"exercise_date"`
}

func (td TodayExercise) todayExerciseInsertQuery(db *sql.DB) (int, error) {
	var seq int
	stmt, err := db.Prepare(
		`INSERT INTO 
			t_today_exercises(trainer_id, 
				group_name, 
				exercises,
				created_user,
				updated_user, 
				user_id, 
				exercise_date) 
			VALUES(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		seq = 0
		return seq, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(
		td.Trainer_Id, td.Group_Name, td.Exercises, td.Trainer_Id, td.Trainer_Id, td.User_Id, td.Exercise_Date)
	if err != nil {
		seq = 0
		return seq, errors.Wrap(err, "Failed to insert to Database")
	}
	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		seq = 0
		return seq, errors.Wrap(err, "Failed to insert last id to Database")
	}
	seq = int(id)

	return seq, nil
}

func (td TodayExercise) selectTodayExercises(db *sql.DB, start_date, end_date string) ([]TodayExercise, error) {
	var tdExercises []TodayExercise
	nullTdExercises := [0]TodayExercise{}
	rows, err := db.Query(
		`SELECT
			t.seq,
			t.trainer_id, 
			t.group_name, 
			t.exercises, 
			t.created_date, 
			t.updated_date,
			t.created_user,
			t.updated_user,
			t.user_id,
			u.name,
			t.exercise_date FROM t_today_exercises t left join t_user u on t.user_id = u.id
		WHERE t.trainer_id = ? AND t.group_name = ? AND t.exercise_date >= ? AND t.exercise_date <= ?`, td.Trainer_Id, td.Group_Name, start_date, end_date)
	if err != nil {
		return nullTdExercises[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var tdExercise TodayExercise
		rows.Scan(
			&tdExercise.Seq,
			&tdExercise.Trainer_Id,
			&tdExercise.Group_Name,
			&tdExercise.Exercises,
			&tdExercise.Created_Date,
			&tdExercise.Updated_Date,
			&tdExercise.Created_User,
			&tdExercise.Updated_User,
			&tdExercise.User_Id,
			&tdExercise.User_Name,
			&tdExercise.Exercise_Date)
		tdExercises = append(tdExercises, tdExercise)
	}
	defer rows.Close()

	return tdExercises, nil
}

func (td TodayExercise) deleteTodayExercises(db *sql.DB) (int, error) {
	var rows int
	stmt, err := db.Prepare(
		`DELETE FROM t_today_exercises WHERE seq = ?`)
	if err != nil {
		rows = 0
		return rows, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(td.Seq)
	if err != nil {
		rows = 0
		return rows, errors.Wrap(err, "Failed to delete category")
	}

	row, err := result.RowsAffected()
	if err != nil {
		rows = 0
		return rows, errors.Wrap(err, "Failed to recive From Database")
	}
	defer stmt.Close()
	rows = int(row)

	return rows, err

}

func (td TodayExercise) modifyTodayExercises(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare(
		`UPDATE t_today_exercises 
			SET
				user_id = ?,
				exercises = ?,
				updated_date = now(), 
				updated_user = ?,
				exercise_date = ?
			WHERE seq = ?`)
	if err != nil {
		rows = 0
		return rows, err
	}

	result, err := stmt.Exec(td.User_Id, td.Exercises, td.Trainer_Id, td.Exercise_Date, td.Seq)
	if err != nil {
		rows = 0
		return rows, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		rows = 0
		return rows, err
	}
	defer stmt.Close()
	rows = int(row)

	return
}

func PatchTodayExercises(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		tdExercises := TodayExercise{}
		//trainer_id, group_name := common.GetQueryString(c)
		//exercise_date := c.Query("exercise_date")
		//user_id := c.Query("user_id")
		if err := c.Bind(&tdExercises); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, tdExercises))
			return
		}

		rows, err := tdExercises.modifyTodayExercises(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, rows))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(rows))
		}
	}

	return resultFunc
}

func DeleteTodayExercises(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("t_seq")
		Seq, err := strconv.Atoi(seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, seq))
			return
		}

		tdExercises := TodayExercise{}
		tdExercises.Seq = Seq
		row, err := tdExercises.deleteTodayExercises(db)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "no rows in result set"):
				c.JSON(http.StatusBadRequest, common.FailedResponse(err, row))
			default:
				c.JSON(http.StatusInternalServerError, common.FailedResponse(err, row))
			}
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(row))
		}
	}

	return resultFunc
}

func PostTodayExercises(db *sql.DB) gin.HandlerFunc {
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

func GetTodayExercises(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		start_date := c.Query("start_date")
		end_date := c.Query("end_date")
		tdExercise := TodayExercise{}
		tdExercise.Trainer_Id, tdExercise.Group_Name = trainer_id, group_name
		tdExercises, err := tdExercise.selectTodayExercises(db, start_date, end_date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, tdExercises))
			return
		}

		c.JSON(http.StatusOK, common.SucceedResponse(tdExercises))
	}

	return resultFunc
}

func TodayExerciseRouter(router *gin.Engine, db *sql.DB) {
	tdExercises := router.Group("/api/t/exercises")

	// Get All TodayExercises with specific trainer_id and group_name
	// curl http://127.0.0.1:8080/api/t/exercises?trainer_id=Park&group_name=dygym -X GET
	tdExercises.GET("/", GetTodayExercises(db))

	// curl http://127.0.0.1:8080/api/t/exercises?trainer_id=Park&group_name=dygym -X POST -d '{"exercises": "스쿼트 20회 5Set", "created_user": "Lee", "updated_user": "Lee", "user_id": "Customer2", "exercise_date": "2022-01-18"}' -H "Content-Type: application/json"
	tdExercises.POST("/", PostTodayExercises(db))

	tdExercises.DELETE("/:t_seq", DeleteTodayExercises(db))

	tdExercises.PATCH("/", PatchTodayExercises(db))
}

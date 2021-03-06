package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type Exercise struct {
	Seq          int     `json:"seq"`
	Title        string  `json:"title" binding:"required"`
	Desc         *string `json:"desc"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Category_Seq int     `json:"category_seq"`
}

// This function is that Query all exercise rows
func (e Exercise) selectAllExercise(db *sql.DB, trainer_id, group_name string) ([]Exercise, error) {
	var exercises []Exercise
	nullExercises := [0]Exercise{}
	rows, err := db.Query(
		`SELECT 
			seq, 
			title,
			`+"`desc`,"+
			`group_name,
			trainer_id,
			created_date, 
			updated_date,
			created_user,
			updated_user, 
			category_seq FROM t_exercise 
		WHERE trainer_id = ? AND group_name = ?`, trainer_id, group_name)
	if err != nil {
		return nullExercises[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var exercise Exercise
		rows.Scan(
			&exercise.Seq,
			&exercise.Title,
			&exercise.Desc,
			&exercise.Group_Name,
			&exercise.Trainer_Id,
			&exercise.Created_Date,
			&exercise.Updated_Date,
			&exercise.Created_User,
			&exercise.Updated_User,
			&exercise.Category_Seq)
		exercises = append(exercises, exercise)
	}
	defer rows.Close()

	if len(exercises) > 0 {
		return exercises, nil
	} else {
		return nullExercises[:], errors.Wrap(err, "Exercise less than 1")
	}
}

func (e Exercise) addExercise(db *sql.DB) (int, error) {
	var id int
	stmt, err := db.Prepare(
		`INSERT INTO 
			t_exercise(title, 
				` + "`desc`," + `
				group_name, 
				trainer_id, 
				created_user, 
				updated_user, 
				category_seq) 
			VALUE(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return id, errors.Wrap(err, "Failed to prepare to Database")
	}

	result, err := stmt.Exec(
		e.Title, e.Desc, e.Group_Name, e.Trainer_Id, e.Trainer_Id, e.Trainer_Id, e.Category_Seq)
	if err != nil {
		return id, errors.Wrap(err, "Failed to execute to Database")
	}

	seq, err := result.LastInsertId()
	if err != nil {
		return id, errors.Wrap(err, "Failed to insert last id to Database")
	}

	id = int(seq)
	defer stmt.Close()

	return id, nil
}

func (e Exercise) deleteExercise(db *sql.DB) (int, error) {
	var rows int
	stmt, err := db.Prepare(
		`DELETE FROM t_exercise WHERE seq = ?`)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(e.Seq)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to delete category")
	}

	row, err := result.RowsAffected()
	if err != nil {
		return rows, errors.Wrap(err, "Failed to recive From Database")
	}
	defer stmt.Close()
	rows = int(row)

	return rows, nil
}

func (e Exercise) modifyExercise(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare(
		`UPDATE t_exercise 
			SET 
				title = ?, 
				` + "`desc`" + `= ?, 
				updated_date = now(), 
				updated_user = ? 
			WHERE seq = ?`)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(e.Title, e.Desc, e.Trainer_Id, e.Seq)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to execute to Database")
	}

	row, err := result.RowsAffected()
	if err != nil {
		return rows, errors.Wrap(err, "Failed to receive From Database")
	}
	defer stmt.Close()
	rows = int(row)

	return
}

func GetAllExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		exercise := Exercise{}
		exercises, err := exercise.selectAllExercise(db, trainer_id, group_name)
		nullExercise := [0]Exercise{}
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, nullExercise))
		} else {
			switch {
			default:
				c.JSON(http.StatusOK, common.SucceedResponse(exercises))
			case len(exercises) == 0:
				c.JSON(http.StatusOK, common.SucceedResponse(nullExercise))
			}

		}
	}

	return resultFunc
}

func PostExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		//category_seq := c.Param("category_seq")
		exercise := Exercise{}
		if err := c.Bind(&exercise); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, exercise))

			return
		}

		_, err := common.IsDuplicatedTitle("t_exercise", exercise.Title, db)
		switch err {
		case nil:
			_, err := exercise.addExercise(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercise))
			} else {
				c.JSON(http.StatusCreated, common.SucceedResponse(exercise))
			}

		case common.DuplicatedTitle:
			c.JSON(http.StatusConflict, common.FailedResponse(err, exercise.Title))

		default:
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercise))
		}

	}

	return resultFunc
}

func DeleteExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("exercise_seq")
		Seq, err := strconv.Atoi(seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, seq))
			return
		}

		exercise := Exercise{Seq: Seq}
		rows, err := exercise.deleteExercise(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, rows))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(rows))
		}
	}

	return resultFunc
}

func PatchExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		exercise := Exercise{}
		err := c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, exercise))
			return
		}

		row, err := common.IsDuplicatedTitle("t_exercise", exercise.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercise.Title))
		} else {
			switch row {
			case 0:
				rows, err := exercise.modifyExercise(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, common.FailedResponse(err, rows))
				} else {
					c.JSON(http.StatusCreated, common.SucceedResponse(rows))
				}

			case 1:
				c.JSON(http.StatusConflict, common.SucceedResponse(exercise))
			}
		}

	}

	return resultFunc
}

func ExerciseRouter(router *gin.Engine, db *sql.DB) {
	exercise := router.Group("/api/exercise")

	// GET All Exercise.
	// curl http://127.0.0.1:8080/api/exercise/all?trainer_id=Lee&group_name=kygym -X GET
	exercise.GET("/all", GetAllExercise(db))

	// Create Exercise.
	// curl http://127.0.0.1:8080/api/exercise -X POST -d '{{"title": "?????????5", "desc": "?????? ????????? ?????? 409 ?????? ?????????", "group_name": "dygym", "trainer_id": "trainer1", "category_seq": 23}}' -H "Content-Type: application/json"
	exercise.POST("/", PostExercise(db))

	// Delete Exercise.
	// curl http://127.0.0.1:8080/api/exercise/2?trainer_id=Lee&group_name=kygym -X DELETE
	exercise.DELETE("/:exercise_seq", DeleteExercise(db))

	// Update Exercise.
	// curl http://127.0.0.1:8080/api/exercise/16 -X PATCH -d {"title": "????????? ????????????", "desc": "Blah Blah.."}
	exercise.PATCH("/", PatchExercise(db))
}

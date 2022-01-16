package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func (e Exercise) selectAllExercise(db *sql.DB) (exercises []Exercise, err error) {
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
		WHERE trainer_id = ? AND group_name = ?`, e.Trainer_Id, e.Group_Name)
	if err != nil {
		return
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

	return
}

func (e Exercise) insertExercise(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare(
		`INSERT INTO 
			t_exercise(title, 
				` + "`desc`," + `
				group_name, 
				trainer_id, 
				created_user, 
				updated_user, 
				category_seq) 
			VALUE(?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}

	result, err := stmt.Exec(
		e.Title, e.Desc, e.Group_Name, e.Trainer_Id, e.Created_User, e.Updated_User, e.Category_Seq)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	Id = int(id)
	defer stmt.Close()

	return
}

func (e Exercise) deleteExercise(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare(
		`DELETE FROM t_exercise 
		WHERE seq = ? AND trainer_id = ? AND group_name = ?`)
	if err != nil {
		return
	}

	result, err := stmt.Exec(e.Seq, e.Trainer_Id, e.Group_Name)
	if err != nil {
		return
	}

	row, err := result.RowsAffected()
	if err != nil {
		return
	}
	defer stmt.Close()
	rows = int(row)
	return
}

func (e Exercise) updateExercise(db *sql.DB) (rows int, err error) {
	// Case 1 -> Only Change Title, Case 2 -> Only Change Description, Case 3 -> Change Title and Description.
	if len(e.Title) == 0 || e.Title == "" {
		stmt, err := db.Prepare(
			`UPDATE t_exercise 
				SET 
					` + "`desc`" + `= ?, 
					pdated_date = now(), 
					updated_user = ? 
				WHERE seq = ?`)
		if err != nil {
			rows = 0
			return rows, err
		}

		result, err := stmt.Exec(e.Desc, e.Updated_User, e.Seq)
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
	} else {
		switch {
		case e.Desc == nil:
			stmt, err := db.Prepare(
				`UPDATE t_exercise 
					SET 
						title = ?, 
						updated_date = now(), 
						updated_user = ? 
					WHERE seq = ?`)
			if err != nil {
				rows = 0
				return rows, err
			}

			result, err := stmt.Exec(e.Title, e.Updated_User, e.Seq)
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
		default:
			stmt, err := db.Prepare(
				`UPDATE t_exercise 
					SET 
						title = ?, 
						` + "`desc`" + `= ?, 
						updated_date = now(), 
						updated_user = ? 
					WHERE seq = ?`)
			if err != nil {
				rows = 0
				return rows, err
			}

			result, err := stmt.Exec(e.Title, e.Desc, e.Updated_User, e.Seq)
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
		}
	}
	return
}

func getExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		exercise := Exercise{}
		exercise.Trainer_Id, exercise.Group_Name = trainer_id, group_name
		exercises, err := exercise.selectAllExercise(db)
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

func postExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category_seq := c.Param("category_seq")
		exercise := Exercise{}
		if err := c.Bind(&exercise); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, category_seq))
			return
		}

		row, err := common.DuplicatedTitleCheck("t_exercise", exercise.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercise.Title))
		} else {
			switch {
			case row == 0:
				exercise.Category_Seq, err = strconv.Atoi(category_seq)
				if err != nil {
					c.JSON(http.StatusBadRequest, common.FailedResponse(err, category_seq))
					return
				}

				Seq, err := exercise.insertExercise(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercise))
				} else {
					c.JSON(http.StatusCreated, common.SucceedResponse(Seq))
				}
			default:
				c.JSON(http.StatusCreated, common.SucceedResponse(row))
			}
		}
	}

	return resultFunc
}

func deleteExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("exercise_seq")
		Seq, err := strconv.Atoi(seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, seq))
			return
		}

		trainer_id, group_name := common.GetQueryString(c)
		exercise := Exercise{}
		exercise.Seq, exercise.Trainer_Id, exercise.Group_Name = Seq, trainer_id, group_name
		rows, err := exercise.deleteExercise(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, rows))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(rows))
		}
	}

	return resultFunc
}

func patchExercise(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("exercise_seq")
		Seq, err := strconv.Atoi(seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, seq))
			return
		}

		trainer_id, group_name := common.GetQueryString(c)
		exercise := Exercise{}
		exercise.Seq, exercise.Trainer_Id, exercise.Group_Name = Seq, trainer_id, group_name
		err = c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, exercise))
			return
		}

		rows, err := exercise.updateExercise(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, rows))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(rows))
		}
	}

	return resultFunc
}

func ExerciseRouter(router *gin.Engine, db *sql.DB) {
	exercise := router.Group("/api/exercise")

	// GET All Exercise.
	// curl http://127.0.0.1:8080/api/exercise/all?trainer_id=Lee&group_name=kygym -X GET
	exercise.GET("/all", getExercise(db))

	// Create Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/2 -X POST -d '{"title": "벤치 프레스", "desc": "가슴 근육 향상"}' -H "Content-Type: application/json"
	exercise.POST("/:category_seq", postExercise(db))

	// Delete Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/2?trainer_id=Lee&group_name=kygym -X DELETE
	exercise.DELETE("/:exercise_seq", deleteExercise(db))

	// Update Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/16 -X PATCH -d {"title": "변경된 카테고리", "desc": "Blah Blah.."}
	exercise.PATCH("/:exercise_seq", patchExercise(db))
}

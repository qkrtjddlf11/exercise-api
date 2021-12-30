package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type Exercise struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Desc         *string `json:"desc"`
	User_Id      int     `json:"user_id"`
	Group_Id     int     `json:"group_id"`
	Trainer_Id   int     `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Category_Id  int     `json:"category_id"`
}

// This function is that Query all exercise rows
func (e Exercise) exerciseGetQueryAll(db *sql.DB) (exercises []Exercise, err error) {
	rows, err := db.Query(
		`SELECT id, 
		title,
		` + "`desc`," +
			`user_id,
		group_id,
		trainer_id,
		created_date, 
		updated_date,
		created_user,
		updated_user, 
		category_id FROM t_exercise`)
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise Exercise
		rows.Scan(
			&exercise.Id,
			&exercise.Title,
			&exercise.Desc,
			&exercise.User_Id,
			&exercise.Group_Id,
			&exercise.Trainer_Id,
			&exercise.Created_Date,
			&exercise.Updated_Date,
			&exercise.Created_User,
			&exercise.Updated_User,
			&exercise.Category_Id)
		exercises = append(exercises, exercise)
	}
	defer rows.Close()

	return
}

func (e Exercise) exerciseInsertQuery(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_exercise(title, `desc`, user_id, group_id, trainer_id, created_user, updated_user, category_id) VALUE(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	result, err := stmt.Exec(
		e.Title, e.Desc, e.User_Id, e.Group_Id, e.Trainer_Id, e.Created_User, e.Updated_User, e.Category_Id)
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

func (e Exercise) exerciseDeleteQuery(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare(
		"DELETE FROM t_exercise WHERE id = ?")
	if err != nil {
		return
	}

	result, err := stmt.Exec(e.Id)
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

func (e Exercise) exerciseUpdateQuery(db *sql.DB) (rows int, err error) {
	// Case 1 -> Only Change Title, Case 2 -> Only Change Description, Case 3 -> Change Title and Description.
	if len(e.Title) == 0 || e.Title == "" {
		stmt, err := db.Prepare(
			"UPDATE t_exercise SET `desc` = ?, updated_date = now(), updated_user = ? WHERE id = ?")
		if err != nil {
			rows = 0
			return rows, err
		}

		result, err := stmt.Exec(e.Desc, e.Updated_User, e.Id)
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
				"UPDATE t_exercise SET title = ?, updated_date = now(), updated_user = ? WHERE id = ?")
			if err != nil {
				rows = 0
				return rows, err
			}

			result, err := stmt.Exec(e.Title, e.Updated_User, e.Id)
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
				"UPDATE t_exercise SET title = ?, `desc` = ?, updated_date = now(), updated_user = ? WHERE id = ?")
			if err != nil {
				rows = 0
				return rows, err
			}

			result, err := stmt.Exec(e.Title, e.Desc, e.Updated_User, e.Id)
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

func ExerciseRouter(router *gin.Engine, db *sql.DB) {
	exercise := router.Group("/api/exercise")

	// GET All Exercise.
	// curl http://127.0.0.1:8080/api/exercise/all -X GET
	exercise.GET("/all", func(c *gin.Context) {
		exercise := Exercise{}
		exercises, err := exercise.exerciseGetQueryAll(db)
		if err != nil {
			nullExercise := [0]Exercise{}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"value":   nullExercise,
			})
		} else {
			c.JSON(http.StatusOK, exercises)
		}
	})

	// Create Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/2 -X POST -d '{"title": "벤치 프레스", "desc": "가슴 근육 향상"}' -H "Content-Type: application/json"
	exercise.POST("/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		exercise := Exercise{}
		err := c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		row, err := common.DuplicatedTitleCheck("t_exercise", exercise.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			switch {
			case row == 0:
				exercise.Category_Id, _ = strconv.Atoi(category_id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": err.Error(),
					})
					return
				}

				Id, err := exercise.exerciseInsertQuery(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": fmt.Sprintf(" %d Successfully Created", Id),
					})
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Duplicated Title"),
				})
			}
		}
	})

	// Delete Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/5 -X DELETE
	exercise.DELETE("/:exercise_id", func(c *gin.Context) {
		id := c.Param("exercise_id")
		Id, err := strconv.ParseInt(id, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		exercise := Exercise{Id: int(Id)}
		rows, err := exercise.exerciseDeleteQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			if rows > 0 {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully Deleted exercise_id: %d", Id),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing Deleted exercise_id : %d", Id),
				})
			}
		}
	})

	// Update Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/16 -X PATCH -d {"title": "변경된 카테고리", "desc": "Blah Blah.."}
	exercise.PATCH("/:exercise_id", func(c *gin.Context) {
		id := c.Param("exercise_id")
		Id, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		exercise := Exercise{}
		exercise.Id = Id
		err = c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		rows, err := exercise.exerciseUpdateQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			if rows > 0 {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully Updated exercise_id : %d", Id),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing Updated, Check exercise_id : %d", Id),
				})
			}
		}
	})
}

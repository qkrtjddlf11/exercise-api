package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Exercise struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Desc        string  `json:"desc"`
	Created_At  *string `json:"createdAt"`
	Updated_At  *string `json:"updatedAt"`
	Category_Id int     `json:"category_id"`
}

func errCheck(err error) error {
	if err != nil {
		return err
	}

	return nil
}

// This function is to Check Duplicated Name in category or exercise table' name
func duplicatedNameCheck(table, name string, db *sql.DB) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(name) FROM %s WHERE name = '%s'", table, name)
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 1, err
	}

	return count, nil
}

// This function is that Query all exercise rows
func (e Exercise) exerciseGetQueryAll(db *sql.DB) (exercises []Exercise, err error) {
	rows, err := db.Query("SELECT id, name, `desc`, createdAt, updatedAt, category_id FROM exercise")
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise Exercise
		rows.Scan(&exercise.Id, &exercise.Name, &exercise.Desc, &exercise.Created_At, &exercise.Updated_At, &exercise.Category_Id)
		exercises = append(exercises, exercise)
	}
	defer rows.Close()

	return
}

func (e Exercise) exerciseInsertQuery(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare("INSERT INTO exercise(name, `desc`, category_id) VALUE(?, ?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(e.Name, e.Desc, e.Category_Id)
	if err != nil {
		return
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return
	}

	Id = int(id)
	defer stmt.Close()

	return
}

func (e Exercise) exerciseDeleteQuery(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare("DELETE FROM exercise WHERE id = ?")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(e.Id)
	if err != nil {
		return
	}

	row, err := rs.RowsAffected()
	if err != nil {
		return
	}
	defer stmt.Close()
	rows = int(row)
	return
}

func (e Exercise) exerciseUpdateQuery(db *sql.DB) (rows int, err error) {
	if len(e.Name) == 0 || e.Name == "" {
		stmt, err := db.Prepare("UPDATE exercise SET `desc` = ?, updatedAt = now() WHERE id = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, e.Id)
		errCheck(err)

		row, err := rs.RowsAffected()

		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE exercise SET name = ?, `desc` = ?, updatedAt = now() WHERE id = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, e.Desc, e.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
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
			c.JSON(http.StatusInternalServerError, nullExercise)
		} else {
			c.JSON(http.StatusOK, exercises)
		}
	})

	// Create Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/create/2 -X POST -d '{"name": "벤치 프레스", "desc": "가슴 근육 향상"}' -H "Content-Type: application/json"
	exercise.POST("/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		exercise := Exercise{}
		err := c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
			return
		}

		row, err := duplicatedNameCheck("exercise", exercise.Name, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed create exercise"),
			})
		} else {
			switch {
			case row == 0:
				exercise.Category_Id, _ = strconv.Atoi(category_id)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": fmt.Sprintf("Invalid Parameter"),
					})
					return
				}

				Id, err := exercise.exerciseInsertQuery(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": fmt.Sprintf("Failed Create"),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": fmt.Sprintf(" %d Successfully Created", Id),
					})
				}
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("Duplicated Name"),
				})
			}
		}
	})

	// Delete Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/delete/5 -X DELETE
	exercise.DELETE("/:exercise_id", func(c *gin.Context) {
		id := c.Param("exercise_id")
		Id, err := strconv.ParseInt(id, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		exercise := Exercise{Id: int(Id)}
		rows, err := exercise.exerciseDeleteQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed delete exercise"),
			})
		} else {
			if rows > 0 {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully deleted exercise_id: %d", Id),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing deleted exercise_id : %d", Id),
				})
			}
		}
	})

	// Update Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/patch/16 -X PATCH -d {"name": "변경된 카테고리", "desc": "Blah Blah.."}
	exercise.PATCH("/:exercise_id", func(c *gin.Context) {
		id := c.Param("exercise_id")
		Id, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		exercise := Exercise{}
		exercise.Id = Id
		err = c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
			return
		}

		rows, err := exercise.exerciseUpdateQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed update exercise_id : %d", Id),
			})
		} else {
			if rows > 0 {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully update exercise_id : %d", Id),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing updated exercise_id : %d", Id),
				})
			}

		}
	})
}

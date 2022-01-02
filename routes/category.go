package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type Category struct {
	Seq          int     `json:"seq"`
	Title        string  `json:"title" binding:"required"`
	Desc         *string `json:"desc" binding:"required"`
	Group_Id     string  `json:"group_id" binding:"required"`
	Trainer_Id   string  `json:"trainer_id" binding:"required"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user" binding:"required"`
	Updated_User string  `json:"updated_user" binding:"required"`
	Count        int     `json:"count"`
}

type ExerciseInCatetory struct {
	Seq   int    `json:"seq"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// This function is that Query all category rows
func (c Category) selectAllCategory(db *sql.DB) (categories []Category, err error) {
	rows, err := db.Query(
		`SELECT c.seq,
		c.title,
		` + "c.`desc`," +
			`c.group_id,
		c.trainer_id,
		c.created_date,
		c.updated_date,
		c.created_user,
		c.updated_user,
		COUNT(e.category_seq) AS count FROM t_category c left join t_exercise e on e.category_seq = c.seq group by c.seq`)
	if err != nil {
		return
	}

	for rows.Next() {
		var category Category
		rows.Scan(
			&category.Seq,
			&category.Title,
			&category.Desc,
			&category.Group_Id,
			&category.Trainer_Id,
			&category.Created_Date,
			&category.Updated_Date,
			&category.Created_User,
			&category.Updated_User,
			&category.Count)
		categories = append(categories, category)
	}
	defer rows.Close()

	return
}

// This function is Insert category
func (c Category) insertCategory(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_category(title, `desc`, group_id, trainer_id, created_user, updated_user) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	result, err := stmt.Exec(
		c.Title, c.Desc, c.Group_Id, c.Trainer_Id, c.Created_User, c.Updated_User)
	if err != nil {
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

// This function is Select all exercises has category_seq
func (e ExerciseInCatetory) selectExerciseInCategory(category_seq int, db *sql.DB) (exercies []ExerciseInCatetory, err error) {
	rows, err := db.Query(
		"SELECT seq, title, `desc` FROM t_exercise WHERE category_seq = ?", category_seq)
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise ExerciseInCatetory
		rows.Scan(&exercise.Seq, &exercise.Title, &exercise.Desc)
		exercies = append(exercies, exercise)
	}
	defer rows.Close()

	return
}

// This function is Delete category
func (c Category) deleteCategory(db *sql.DB) (rows int, err error) {
	var inCategory int
	count := db.QueryRow(
		"SELECT COUNT(e.category_seq) AS count FROM t_category c left join t_exercise e on e.category_seq = c.seq WHERE c.seq = ? group by c.seq", c.Seq)
	err = count.Scan(&inCategory)
	if err != nil {
		rows = 0
		return
	}

	if inCategory > 0 {
		rows = 0
		return
	}

	stmt, err := db.Prepare(
		"DELETE FROM t_category WHERE seq = ?")
	if err != nil {
		rows = 0
		return
	}

	result, err := stmt.Exec(c.Seq)
	if err != nil {
		rows = 0
		return
	}

	row, err := result.RowsAffected()
	if err != nil {
		rows = 0
		return
	}
	defer stmt.Close()
	rows = int(row)

	return
}

// This function is Update category
func (c Category) updateCategory(db *sql.DB) (rows int, err error) {
	// Case 1 -> Only Change Title, Case 2 -> Only Change Description, Case 3 -> Change Title and Description.
	if c.Title == "" {
		stmt, err := db.Prepare(
			"UPDATE t_category SET `desc` = ?, updated_date = now(), updated_user = ? WHERE seq = ?")
		if err != nil {
			rows = 0
			return rows, err
		}

		result, err := stmt.Exec(c.Desc, c.Updated_User, c.Seq)
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
		case c.Desc == nil:
			stmt, err := db.Prepare(
				"UPDATE t_category SET title = ?, updated_date = now(), updated_user = ? WHERE seq = ?")
			if err != nil {
				rows = 0
				return rows, err
			}
			result, err := stmt.Exec(c.Title, c.Updated_User, c.Seq)
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
				"UPDATE t_category SET title = ?, `desc` = ?, updated_date = now(), updated_user = ? WHERE seq = ?")
			if err != nil {
				rows = 0
				return rows, err
			}
			result, err := stmt.Exec(c.Title, c.Desc, c.Updated_User, c.Seq)
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

func getAllCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category := Category{}
		nullCategory := [0]Category{}
		categories, err := category.selectAllCategory(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"paramegers": gin.H{
					"body": nullCategory,
				},
			})
		} else {
			switch {
			case len(categories) > 0:
				c.JSON(http.StatusOK, categories)
			default:
				c.JSON(http.StatusOK, nullCategory)
			}
		}
	}

	return resultFunc
}

func getExercisesInCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category_seq := c.Param("category_seq")
		Category_Seq, err := strconv.Atoi(category_seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"params": gin.H{
						"category_seq": category_seq,
					},
				},
			})
			return
		}

		nullExercise := [0]ExerciseInCatetory{}
		exercise := ExerciseInCatetory{}
		exercises, err := exercise.selectExerciseInCategory(Category_Seq, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"paramegers": gin.H{
					"body": nullExercise,
				},
			})
		} else {
			switch {
			case len(exercises) > 0:
				c.JSON(http.StatusOK, exercises)
			default:
				c.JSON(http.StatusOK, nullExercise) // return []
			}
		}
	}

	return resultFunc
}

func postCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category := Category{}
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"paramters": gin.H{
					"body": category,
				},
			})
			return
		}

		row, err := common.DuplicatedTitleCheck("t_category", category.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			switch {
			case row == 0:
				_, err := category.insertCategory(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
						"parameters": gin.H{
							"body": category,
						},
					})
				} else {
					c.JSON(http.StatusOK, gin.H{})
				}
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Duplicated Title"),
					"parameters": gin.H{
						"body": category,
					},
				})

			}
		}
	}

	return resultFunc
}

func deleteCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("category_seq")
		Seq, err := strconv.ParseInt(seq, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"params": gin.H{
						"category_seq": Seq,
					},
				},
			})
			return
		}

		category := Category{Seq: int(Seq)}
		row, err := category.deleteCategory(db)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "no rows in result set"):
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
					"parameters": gin.H{
						"parameter": gin.H{
							"category_seq": Seq,
						},
					},
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}

		} else {
			switch {
			case row > 0:
				c.JSON(http.StatusOK, gin.H{})
			}
		}
	}

	return resultFunc
}

func patchCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		seq := c.Param("category_seq")
		Seq, err := strconv.Atoi(seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"parameter": gin.H{
						"category_seq": Seq,
					},
				},
			})
			return
		}

		category := Category{}
		category.Seq = Seq
		if err = c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"parameter": gin.H{
						"category_seq": Seq,
						"body":         category,
					},
				},
			})
			return
		}

		row, err := category.updateCategory(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"value":   row,
			})
		} else {
			switch {
			case row > 0:
				c.JSON(http.StatusOK, gin.H{})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Nothing Updated, Check category_seq: %d", Seq),
				})
			}
		}
	}

	return resultFunc
}

func CategoryRouter(router *gin.Engine, db *sql.DB) {
	category := router.Group("/api/category")
	// GET All Category.
	// curl http://127.0.0.1:8080/api/category/all -X GET
	category.GET("/all", getAllCategory(db))

	// GET All Exercises in Specific Category.
	// curl http://127.0.0.1:8080/api/category/exercise/5 -X GET
	category.GET("/exercise/:category_seq", getExercisesInCategory(db))

	// Create Category
	// curl http://127.0.0.1:8080/api/category -X POST -d '{"title": "등", "desc": "등 근육의 전반적인 향상", "group_id": 1, "trainer_id": "Choi Trainer","created_user": "Park", "updated_user": "Park"}' -H "Content-Type: application/json"
	category.POST("/", postCategory(db))

	// Delete Specific Category.
	// curl http://127.0.0.1:8080/api/category/delete/5 -X DELETE
	category.DELETE("/:category_seq", deleteCategory(db))

	// Update Specific Category.
	// curl http://127.0.0.1:8080/api/category/patch/16 -X PATCH -d {"title": "변경된 카테고리", "desc": "Blah Blah.."}
	category.PATCH("/:category_seq", patchCategory(db))
}

package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type Category struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Desc         *string `json:"desc"`
	Group_Id     int     `json:"group_id"`
	Trainer_Id   int     `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Count        int     `json:"count"`
}

type ExerciseInCatetory struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

// This function is that Query all category rows
func (c Category) categoryGetQueryAll(db *sql.DB) (categories []Category, err error) {
	rows, err := db.Query(
		`SELECT c.id, 
		c.title, 
		` + "c.`desc`," +
			`c.group_id, 
		c.trainer_id, 
		c.created_date, 
		c.updated_date, 
		c.created_user, 
		c.updated_user, 
		COUNT(e.category_id) AS count FROM t_category c left join t_exercise e on e.category_id = c.id group by c.id`)
	if err != nil {
		return
	}

	for rows.Next() {
		var category Category
		rows.Scan(
			&category.Id,
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
func (c Category) categoryInsertQuery(db *sql.DB) (Id int, err error) {
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

// This function is Select all exercises has category_id
func (e ExerciseInCatetory) exerciseGetQueryInCategory(category_id int, db *sql.DB) (exercies []ExerciseInCatetory, err error) {
	rows, err := db.Query(
		"SELECT id, title, `desc` FROM t_exercise WHERE category_id = ?", category_id)
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise ExerciseInCatetory
		rows.Scan(&exercise.Id, &exercise.Title, &exercise.Desc)
		exercies = append(exercies, exercise)
	}
	defer rows.Close()

	return
}

// This function is Delete category
func (c Category) categoryDeleteQuery(db *sql.DB) (rows int, err error) {
	var inCategory int
	count := db.QueryRow(
		"SELECT COUNT(e.category_id) AS count FROM t_category c left join t_exercise e on e.category_id = c.id group by c.id")
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
		"DELETE FROM t_category WHERE id = ?")
	if err != nil {
		rows = 0
		return
	}

	result, err := stmt.Exec(c.Id)
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
func (c Category) categoryUpdateQuery(db *sql.DB) (rows int, err error) {
	// Case 1 -> Only Change Title, Case 2 -> Only Change Description, Case 3 -> Change Title and Description.
	if c.Title == "" {
		stmt, err := db.Prepare(
			"UPDATE t_category SET `desc` = ?, updated_date = now(), updated_user = ? WHERE id = ?")
		if err != nil {
			rows = 0
			return rows, err
		}

		result, err := stmt.Exec(c.Desc, c.Updated_User, c.Id)
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
				"UPDATE t_category SET title = ?, updated_date = now(), updated_user = ? WHERE id = ?")
			if err != nil {
				rows = 0
				return rows, err
			}
			result, err := stmt.Exec(c.Title, c.Updated_User, c.Id)
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
				"UPDATE t_category SET title = ?, `desc` = ?, updated_date = now(), updated_user = ? WHERE id = ?")
			if err != nil {
				rows = 0
				return rows, err
			}
			result, err := stmt.Exec(c.Title, c.Desc, c.Updated_User, c.Id)
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

func CategoryRouter(router *gin.Engine, db *sql.DB) {
	category := router.Group("/api/category")

	// GET All Category.
	// curl http://127.0.0.1:8080/api/category/all -X GET
	category.GET("/all", func(c *gin.Context) {
		category := Category{}
		categories, err := category.categoryGetQueryAll(db)
		if err != nil {
			nullCategory := [0]Category{}
			c.JSON(http.StatusInternalServerError, nullCategory)
		} else {
			c.JSON(http.StatusOK, categories)
		}
	})

	// GET All Exercises in Specific Category.
	// curl http://127.0.0.1:8080/api/category/exercise/5 -X GET
	category.GET("/exercise/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		Category_id, err := strconv.Atoi(category_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		exercise := ExerciseInCatetory{}
		exercises, err := exercise.exerciseGetQueryInCategory(Category_id, db)
		if err != nil {
			nullExercise := [0]Exercise{}
			c.JSON(http.StatusInternalServerError, nullExercise)
		} else if len(exercises) > 0 {
			c.JSON(http.StatusOK, exercises)
		} else {
			nullExercise := [0]ExerciseInCatetory{}
			c.JSON(http.StatusOK, nullExercise) // return []
		}
	})

	// Create Specific Category.
	// curl http://127.0.0.1:8080/api/category -X POST -d '{"title": "등", "desc": "등 근육의 전반적인 향상", "created_user": "Park"}' -H "Content-Type: application/json"
	category.POST("/", func(c *gin.Context) {
		category := Category{}
		err := c.Bind(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
		}

		row, err := common.DuplicatedTitleCheck("t_category", category.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed Create category"),
			})
		} else {
			switch {
			case row == 0:
				Id, err := category.categoryInsertQuery(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": fmt.Sprintf("ID : %d Successfully Created", Id),
					})
				}
			case row == 1:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Duplicated Title"),
				})
			}
		}

	})

	// Delete Specific Category.
	// curl http://127.0.0.1:8080/api/category/delete/5 -X DELETE
	category.DELETE("/:category_id", func(c *gin.Context) {
		id := c.Param("category_id")
		Id, err := strconv.ParseInt(id, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		category := Category{Id: int(Id)}
		row, err := category.categoryDeleteQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Id)
		} else {
			switch {
			case row > 0:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully Deleted category_id: %d", Id),
				})
			default:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing Deleted category_id: %d", Id),
				})
			}
		}
	})

	// Update Specific Category.
	// curl http://127.0.0.1:8080/api/category/patch/16 -X PATCH -d {"title": "변경된 카테고리", "desc": "Blah Blah.."}
	category.PATCH("/:category_id", func(c *gin.Context) {
		id := c.Param("category_id")
		Id, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		category := Category{}
		category.Id = Id
		err = c.Bind(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
			return
		}

		row, err := category.categoryUpdateQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, row)
		} else {
			switch {
			case row > 0:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully Updated category_id: %d", Id),
				})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Nothing Updated, Check category_id: %d", Id),
				})
			}
		}
	})
}

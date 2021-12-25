package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Category struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Desc       *string `json:"desc"`
	Created_At *string `json:"createdAt"`
	Updated_At *string `json:"updatedAt"`
	Count      int     `json:"count"`
	//Created_User string  `json:"created_user"`
	//Updated_User *string `json:"updated_user"`
}

type ExerciseInCatetory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// This function is that Query all category rows
func (c Category) categoryGetQueryAll(db *sql.DB) (categories []Category, err error) {
	// select c.num, c.name, count(e.category_id) as \`count\` from category c left join exercise e on e.category_id = c.num group by c.num
	rows, err := db.Query("SELECT c.id, c.name, c.`desc`, c.createdAt, c.updatedAt, count(e.category_id) AS count FROM category c left join exercise e on e.category_id = c.id group by c.id")
	if err != nil {
		return
	}

	for rows.Next() {
		var category Category
		rows.Scan(&category.Id, &category.Name, &category.Desc, &category.Created_At, &category.Updated_At, &category.Count)
		categories = append(categories, category)
	}
	defer rows.Close()

	return
}

func (e ExerciseInCatetory) exerciseGetQueryInCategory(category_id int, db *sql.DB) (exercies []ExerciseInCatetory, err error) {
	rows, err := db.Query("SELECT id, name, `desc` FROM exercise WHERE category_id = ?", category_id)
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise ExerciseInCatetory
		rows.Scan(&exercise.Id, &exercise.Name, &exercise.Desc)
		exercies = append(exercies, exercise)
	}
	defer rows.Close()

	return
}

func (c Category) categoryDeleteQuery(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare("DELETE FROM category WHERE id = ?")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(c.Id)
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

func (c Category) categoryUpdateQuery(db *sql.DB) (rows int, err error) {
	if len(c.Name) == 0 || c.Name == "" {
		stmt, err := db.Prepare("UPDATE category SET `desc` = ?, updatedAt = now() WHERE id = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Desc, c.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE category SET name = ?, `desc` = ?, updatedAt = now() WHERE id = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Name, c.Desc, c.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	}

	return
}

func (c Category) categoryInsertQuery(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare("INSERT INTO category(name, `desc`) VALUES(?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(c.Name, c.Desc)
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
	// curl http://127.0.0.1:8080/api/category -X POST -d '{"name": "등", "desc": "등 근육의 전반적인 향상", "created_user": "Park"}' -H "Content-Type: application/json"
	category.POST("/", func(c *gin.Context) {
		category := Category{}
		err := c.Bind(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
		}

		row, err := duplicatedNameCheck("category", category.Name, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed create category"),
			})
		} else {
			switch {
			case row == 0:
				Id, err := category.categoryInsertQuery(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": fmt.Sprintf(" %d Successfully Created", Id),
					})
				}
			case row == 1:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Duplicated Name"),
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
		rows, err := category.categoryDeleteQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Id)
		} else {
			switch {
			case rows > 0:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully deleted category_id: %d", Id),
				})
			default:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing deleted category_id: %d", Id),
				})
			}
		}
	})

	// Update Specific Category.
	// curl http://127.0.0.1:8080/api/category/patch/16 -X PATCH -d {"name": "변경된 카테고리", "desc": "Blah Blah.."}
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

		rows, err := category.categoryUpdateQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, rows)
		} else {
			switch {
			case rows > 0:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Successfully update category_id: %d", Id),
				})
			default:
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("Nothing deleted category_id: %d", Id),
				})
			}
		}
	})
}

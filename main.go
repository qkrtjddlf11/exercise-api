package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const MariaDB string = "mysql"

var db *sql.DB

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

type Exercise struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Desc        string  `json:"desc"`
	Created_At  *string `json:"createdAt"`
	Updated_At  *string `json:"updatedAt"`
	Category_Id int     `json:"category_id"`
}

type ExerciseInCatetory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func errCheck(err error) error {
	if err != nil {
		return err
	}

	return nil
}

// This function is to Check Duplicated Name in category or exercise table' name
func duplicatedNameCheck(table, name string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(name) FROM %s WHERE name = '%s'", table, name)
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 1, err
	}

	return count, nil
}

// This function is that Query all category rows
func (c Category) categoryGetQueryAll() (categories []Category, err error) {
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

// This function is that Query all exercise rows
func (e Exercise) exerciseGetQueryAll() (exercises []Exercise, err error) {
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

func (e ExerciseInCatetory) exerciseGetQueryInCategory(category_id int) (exercies []ExerciseInCatetory, err error) {
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

func (c Category) categoryInsertQuery() (Id int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("INSERT INTO category(name, `desc`, createdAt) VALUES(?, ?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(c.Name, c.Desc, dateTime)
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

func (e Exercise) exerciseInsertQuery() (Id int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO exercise(name, `desc`, createdAt, category_id) VALUE(?, ?, ?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(e.Name, e.Desc, dateTime, e.Category_Id)
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

func (c Category) categoryDeleteQuery() (rows int, err error) {
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

func (e Exercise) exerciseDeleteQuery() (rows int, err error) {
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

func (c Category) categoryUpdateQuery() (rows int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	if len(c.Name) == 0 || c.Name == "" {
		stmt, err := db.Prepare("UPDATE category SET `desc` = ?, updatedAt = ? WHERE id = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Desc, dateTime, c.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE category SET name = ?, `desc` = ?, updatedAt = ? WHERE id = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Name, c.Desc, dateTime, c.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	}

	return
}

func (e Exercise) exerciseUpdateQuery() (rows int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	if len(e.Name) == 0 || e.Name == "" {
		stmt, err := db.Prepare("UPDATE exercise SET `desc` = ?, updatedAt = ? WHERE id = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, dateTime, e.Id)
		errCheck(err)

		row, err := rs.RowsAffected()

		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE exercise SET name = ?, `desc` = ?, updatedAt = ? WHERE id = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, e.Desc, dateTime, e.Id)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	}
	return
}

func main() {
	// Logging to file
	gin.DisableConsoleColor()
	logFile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connection DB
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWD"),
		os.Getenv("DBIPADDR"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"))
	db, err = sql.Open(MariaDB, dbString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	// GET All Category.
	// curl http://127.0.0.1:8080/api/category/all -X GET
	router.GET("/api/category/all", func(c *gin.Context) {
		category := Category{}
		categories, err := category.categoryGetQueryAll()
		if err != nil {
			nullCategory := [0]Category{}
			c.JSON(http.StatusInternalServerError, nullCategory)
		} else {
			c.JSON(http.StatusOK, categories)
		}
	})

	// GET All Exercises in Specific Category.
	// curl http://127.0.0.1:8080/api/category/exercise/5 -X GET
	router.GET("/api/category/exercise/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		Category_id, err := strconv.Atoi(category_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		exercise := ExerciseInCatetory{}
		exercises, err := exercise.exerciseGetQueryInCategory(Category_id)
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
	// curl http://127.0.0.1:8080/api/category/create -X POST -d '{"name": "등", "desc": "등 근육의 전반적인 향상", "created_user": "Park"}' -H "Content-Type: application/json"
	router.POST("/api/category", func(c *gin.Context) {
		category := Category{}
		err := c.Bind(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
		}

		row, err := duplicatedNameCheck("category", category.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed create category"),
			})
		} else {
			switch {
			case row == 0:
				Id, err := category.categoryInsertQuery()
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
	router.DELETE("/api/category/:category_id", func(c *gin.Context) {
		id := c.Param("category_id")
		Id, err := strconv.ParseInt(id, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		category := Category{Id: int(Id)}
		rows, err := category.categoryDeleteQuery()
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
	router.PATCH("/api/category/:category_id", func(c *gin.Context) {
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

		rows, err := category.categoryUpdateQuery()
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

	// GET All Exercise.
	// curl http://127.0.0.1:8080/api/exercise/all -X GET
	router.GET("/api/exercise/all", func(c *gin.Context) {
		exercise := Exercise{}
		exercises, err := exercise.exerciseGetQueryAll()
		if err != nil {
			nullExercise := [0]Exercise{}
			c.JSON(http.StatusInternalServerError, nullExercise)
		} else {
			c.JSON(http.StatusOK, exercises)
		}
	})

	// Create Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/create/2 -X POST -d '{"name": "벤치 프레스", "desc": "가슴 근육 향상"}' -H "Content-Type: application/json"
	router.POST("/api/exercise/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		exercise := Exercise{}
		err := c.Bind(&exercise)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid JSON Format"),
			})
			return
		}

		row, err := duplicatedNameCheck("exercise", exercise.Name)
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

				Id, err := exercise.exerciseInsertQuery()
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
	router.DELETE("/api/exercise/:exercise_id", func(c *gin.Context) {
		id := c.Param("exercise_id")
		Id, err := strconv.ParseInt(id, 10, 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Invalid Parameter"),
			})
			return
		}

		exercise := Exercise{Id: int(Id)}
		rows, err := exercise.exerciseDeleteQuery()
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
	router.PATCH("/api/exercise/:exercise_id", func(c *gin.Context) {
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

		rows, err := exercise.exerciseUpdateQuery()
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

	router.Run("0.0.0.0:8080")
}

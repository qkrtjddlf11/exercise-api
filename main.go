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
	Num         int     `json:"num"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Created_At  *string `json:"createdAt"`
	Updated_At  *string `json:"updatedAt"`
	Count       int     `json:"count"`
	//Created_User string  `json:"created_user"`
	//Updated_User *string `json:"updated_user"`
}

type Exercise struct {
	Num         int     `json:"num"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Created_At  *string `json:"createdAt"`
	Updated_At  *string `json:"updatedAt"`
	Category_Id int     `json:"category_id"`
}

type ExerciseInCatetory struct {
	Num         int    `json:"num"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func errCheck(err error) error {
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c Category) categoryGetQueryAll() (categories []Category, err error) {
	// select c.num, c.name, count(e.category_id) as \`count\` from category c left join exercise e on e.category_id = c.num group by c.num
	rows, err := db.Query("SELECT c.num, c.name, c.description, c.createdAt, c.updatedAt, count(e.category_id) AS count FROM category c left join exercise e on e.category_id = c.num group by c.num")
	if err != nil {
		return
	}

	for rows.Next() {
		var category Category
		rows.Scan(&category.Num, &category.Name, &category.Description, &category.Created_At, &category.Updated_At, &category.Count)
		categories = append(categories, category)
	}
	defer rows.Close()

	return
}

func (e Exercise) exerciseGetQueryAll() (exercises []Exercise, err error) {
	rows, err := db.Query("SELECT num, name, description, createdAt, updatedAt, category_id FROM exercise")
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise Exercise
		rows.Scan(&exercise.Num, &exercise.Name, &exercise.Description, &exercise.Created_At, &exercise.Updated_At, &exercise.Category_Id)
		exercises = append(exercises, exercise)
	}
	defer rows.Close()

	return
}

func (e ExerciseInCatetory) exerciseGetQueryInCategory(category_id int) (exercies []ExerciseInCatetory, err error) {
	rows, err := db.Query("SELECT num, name, description FROM exercise WHERE category_id = ?", category_id)
	if err != nil {
		return
	}

	for rows.Next() {
		var exercise ExerciseInCatetory
		rows.Scan(&exercise.Num, &exercise.Name, &exercise.Description)
		exercies = append(exercies, exercise)
	}
	defer rows.Close()

	return
}

func (c Category) categoryInsertQuery() (Num int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO category(name, description, createdAt) VALUES(?, ?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(c.Name, c.Description, dateTime)
	if err != nil {
		return
	}

	num, err := rs.LastInsertId()
	errCheck(err)

	Num = int(num)
	defer stmt.Close()

	return
}

func (e Exercise) exerciseInsertQuery() (Num int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("INSERT INTO exercise(name, description, createdAt, category_id) VALUE(?, ?, ?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(e.Name, e.Description, dateTime, e.Category_Id)
	if err != nil {
		return
	}

	num, err := rs.LastInsertId()
	errCheck(err)

	Num = int(num)
	defer stmt.Close()

	return
}

func (c Category) categoryDeleteQuery() (rows int, err error) {
	stmt, err := db.Prepare("DELETE FROM category WHERE num = ?")
	errCheck(err)

	rs, err := stmt.Exec(c.Num)
	errCheck(err)

	row, err := rs.RowsAffected()
	errCheck(err)
	defer stmt.Close()
	rows = int(row)

	return
}

func (e Exercise) exerciseDeleteQuery() (rows int, err error) {
	stmt, err := db.Prepare("DELETE FROM exercise WHERE num = ?")
	errCheck(err)

	rs, err := stmt.Exec(e.Num)
	errCheck(err)

	row, err := rs.RowsAffected()
	errCheck(err)
	defer stmt.Close()
	rows = int(row)
	return
}

func (c Category) categoryUpdateQuery() (rows int, err error) {
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	if len(c.Name) == 0 || c.Name == "" {
		stmt, err := db.Prepare("UPDATE category SET description = ?, updatedAt = ? WHERE num = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Description, dateTime, c.Num)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE category SET name = ?, description = ?, updatedAt = ? WHERE num = ?")
		errCheck(err)
		rs, err := stmt.Exec(c.Name, c.Description, dateTime, c.Num)
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
		stmt, err := db.Prepare("UPDATE exercise SET description = ?, updatedAt = ? WHERE num = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, dateTime, e.Num)
		errCheck(err)

		row, err := rs.RowsAffected()
		errCheck(err)
		defer stmt.Close()
		rows = int(row)
	} else {
		stmt, err := db.Prepare("UPDATE exercise SET name = ?, description = ?, updatedAt = ? WHERE num = ?")
		errCheck(err)

		rs, err := stmt.Exec(e.Name, e.Description, dateTime, e.Num)
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
	logFile, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("env [DBNAME]:", os.Getenv("DBNAME"), os.Getenv("DBUSER"))

	// Connection DB
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DBUSER"), os.Getenv("DBPASSWD"), os.Getenv("DBIPADDR"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	db, err = sql.Open(MariaDB, dbString)
	errCheck(err)
	defer db.Close()

	router := gin.Default()

	// GET All Category.
	// curl http://127.0.0.1:8080/api/category/all -X GET
	router.GET("/api/category/all", func(c *gin.Context) {
		category := Category{}
		categories, err := category.categoryGetQueryAll()
		errCheck(err)

		c.JSON(http.StatusOK, categories)
	})

	// GET All Exercises in Specific Category.
	// curl http://127.0.0.1:8080/api/category/exercise/5 -X GET
	router.GET("/api/category/exercise/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		Category_id, err := strconv.Atoi(category_id)
		errCheck(err)

		exercise := ExerciseInCatetory{}
		exercises, err := exercise.exerciseGetQueryInCategory(Category_id)
		errCheck(err)

		if len(exercises) > 0 {
			c.JSON(http.StatusOK, exercises)
		} else {
			nullExercise := [0]ExerciseInCatetory{}
			c.JSON(http.StatusOK, nullExercise) // return []
		}
	})

	// Create Specific Category.
	// curl http://127.0.0.1:8080/api/category/create -X POST -d '{"name": "등", "description": "등 근육의 전반적인 향상", "created_user": "Park"}' -H "Content-Type: application/json"
	router.POST("/api/category", func(c *gin.Context) {
		category := Category{}
		err := c.Bind(&category)
		errCheck(err)

		Num, err := category.categoryInsertQuery()
		errCheck(err)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %d Successfully Created", Num),
		})
	})

	// Delete Specific Category.
	// curl http://127.0.0.1:8080/api/category/delete/5 -X DELETE
	router.DELETE("/api/category/:num", func(c *gin.Context) {
		num := c.Param("num")

		Num, err := strconv.ParseInt(num, 10, 10)
		errCheck(err)

		category := Category{Num: int(Num)}
		rows, err := category.categoryDeleteQuery()
		if err != nil {
			log.Fatal(err, rows)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted num: %d", Num),
		})

	})

	// Update Specific Category.
	// curl http://127.0.0.1:8080/api/category/patch/16 -X PATCH -d {"name": "변경된 카테고리", "description": "Blah Blah.."}
	router.PATCH("/api/category/:num", func(c *gin.Context) {
		num := c.Param("num")

		Num, err := strconv.Atoi(num)
		errCheck(err)

		category := Category{}
		category.Num = Num
		err = c.Bind(&category)
		errCheck(err)

		rows, err := category.categoryUpdateQuery()
		if err != nil {
			log.Fatal(err, rows)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully update num: %d", Num),
		})
	})

	// GET All Exercise.
	// curl http://127.0.0.1:8080/api/exercise/all -X GET
	router.GET("/api/exercise/all", func(c *gin.Context) {
		exercise := Exercise{}
		exercises, err := exercise.exerciseGetQueryAll()
		errCheck(err)

		c.JSON(http.StatusOK, exercises)
	})

	// Create Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/create/2 -X POST -d '{"name": "벤치 프레스", "description": "가슴 근육 향상"}' -H "Content-Type: application/json"
	router.POST("/api/exercise/:category_id", func(c *gin.Context) {
		category_id := c.Param("category_id")
		exercise := Exercise{}
		err := c.Bind(&exercise)
		errCheck(err)

		exercise.Category_Id, _ = strconv.Atoi(category_id)
		Num, err := exercise.exerciseInsertQuery()
		errCheck(err)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %d Successfully Created", Num),
		})
	})

	// Delete Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/delete/5 -X DELETE
	router.DELETE("/api/exercise/:num", func(c *gin.Context) {
		num := c.Param("num")

		Num, err := strconv.ParseInt(num, 10, 10)
		errCheck(err)

		exercise := Exercise{Num: int(Num)}
		rows, err := exercise.exerciseDeleteQuery()
		if err != nil {
			log.Fatal(err, rows)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted num: %d", Num),
		})

	})

	// Update Specific Exercise.
	// curl http://127.0.0.1:8080/api/exercise/patch/16 -X PATCH -d {"name": "변경된 카테고리", "description": "Blah Blah.."}
	router.PATCH("/api/exercise/:num", func(c *gin.Context) {
		num := c.Param("num")

		Num, err := strconv.Atoi(num)
		errCheck(err)

		exercise := Exercise{}
		exercise.Num = Num
		err = c.Bind(&exercise)
		errCheck(err)

		rows, err := exercise.exerciseUpdateQuery()
		if err != nil {
			log.Fatal(err, rows)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully update num: %d", Num),
		})
	})

	router.Run("0.0.0.0:8080")
}

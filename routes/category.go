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

type Category struct {
	Seq          int     `json:"seq"`
	Title        string  `json:"title"`
	Desc         *string `json:"desc"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Count        int     `json:"count"`
}

type ExerciseInCatetory struct {
	Seq        int    `json:"seq"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Group_Name string `json:"group_name" binding:"required"`
	Trainer_Id string `json:"trainer_id" binding:"required"`
}

// This function is that Query all category rows
func (c Category) selectAllCategory(db *sql.DB, trainer_id, group_name string) ([]Category, error) {
	var categories []Category
	nullCategory := [0]Category{}
	rows, err := db.Query(
		`SELECT 
			c.seq,
			c.title,
			`+"c.`desc`,"+
			`c.group_name,
			c.trainer_id,
			c.created_date,
			c.updated_date,
			c.created_user,
			c.updated_user,
			COUNT(e.category_seq) AS count 
		FROM t_category c left join t_exercise e on e.category_seq = c.seq 
		WHERE c.trainer_id = ? AND c.group_name = ? GROUP BY c.seq`, trainer_id, group_name)
	if err != nil {
		return nullCategory[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var category Category
		rows.Scan(
			&category.Seq,
			&category.Title,
			&category.Desc,
			&category.Group_Name,
			&category.Trainer_Id,
			&category.Created_Date,
			&category.Updated_Date,
			&category.Created_User,
			&category.Updated_User,
			&category.Count)
		categories = append(categories, category)
	}
	defer rows.Close()

	if len(categories) > 0 {
		return categories, nil
	} else {
		return nullCategory[:], errors.Wrap(err, "Category less than 1")
	}
}

/*
func (c Category) preAddCategory(db *sql.DB) error {
	var count int
	row := db.QueryRow(
		`SELECT
			COUNT(title)
		FROM t_category like '%?%'`, c.Title)

	err := row.Scan(&count)
	if err != nil {
		return errors.Wrap(err, "Failed to select From Database")
	}

	if count > 0 {
		return errors.Wrap(err, "Duplicated Category Name")
	}
	return nil
}
*/

// This function is Insert category
func (c Category) addCategory(db *sql.DB) (int, error) {
	var id int
	/*
		err := c.preAddCategory(db)
		if err != nil {
			return id, err
		}
	*/
	stmt, err := db.Prepare(
		`INSERT INTO 
			t_category(title,
				` + "`desc`," + ` 
				group_name, 
				trainer_id, 
				created_user, 
				updated_user) 
			VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return id, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(
		c.Title, c.Desc, c.Group_Name, c.Trainer_Id, c.Trainer_Id, c.Trainer_Id)
	if err != nil {
		return id, errors.Wrap(err, "Failed to insert to Database")
	}
	defer stmt.Close()

	seq, err := result.LastInsertId()
	if err != nil {
		return id, errors.Wrap(err, "Failed to insert last id to Database")
	}
	id = int(seq)

	return id, nil
}

// This function is Select all exercises has category_seq
func (e ExerciseInCatetory) selectExerciseInCategory(category_seq int, db *sql.DB) ([]ExerciseInCatetory, error) {
	var exercies []ExerciseInCatetory
	nullExercises := [0]ExerciseInCatetory{}
	rows, err := db.Query(
		`SELECT 
			seq, 
			title, 
			`+"`desc`,"+` 
			trainer_id, 
			group_name FROM t_exercise 
		WHERE category_seq = ?`, category_seq)
	if err != nil {
		return nullExercises[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var exercise ExerciseInCatetory
		rows.Scan(&exercise.Seq, &exercise.Title, &exercise.Desc, &exercise.Trainer_Id, &exercise.Group_Name)
		exercies = append(exercies, exercise)
	}
	defer rows.Close()

	if len(exercies) > 0 {
		return exercies, nil
	} else {
		return nullExercises[:], errors.Wrap(err, "At least included 1 exercise in Category")
	}
}

// This function is Delete category
func deleteCategory(db *sql.DB, category_seq int) (int, error) {
	var rows int
	var inCategory int
	count := db.QueryRow(
		`SELECT 
			COUNT(e.category_seq) AS count 
		FROM t_category c left join t_exercise e on e.category_seq = c.seq 
		WHERE c.seq = ? GROUP BY c.seq`, category_seq)
	err := count.Scan(&inCategory)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to select before delete Query")
	}

	if inCategory > 0 {
		return rows, errors.Wrap(err, "Included at least 1 exercise in Category")
	}

	stmt, err := db.Prepare(
		`DELETE FROM t_category WHERE seq = ?`)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(category_seq)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to delete category")
	}

	row, err := result.RowsAffected()
	if err != nil {
		return rows, errors.Wrap(err, "Failed to receive From Database")
	}
	defer stmt.Close()
	rows = int(row)

	return rows, nil
}

// This function is Update category
func (c Category) modifyCategory(db *sql.DB) (rows int, err error) {
	stmt, err := db.Prepare(
		`UPDATE t_category 
			SET title = ?,
				` + "`desc`" + `= ?, 
				updated_date = now(), 
				updated_user = ? 
			WHERE seq = ?`)
	if err != nil {
		return rows, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(c.Title, c.Desc, c.Trainer_Id, c.Seq)
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

func GetAllCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		//nullCategory := [0]Category{}
		trainer_id, group_name := common.GetQueryString(c)
		category := Category{}
		categories, err := category.selectAllCategory(db, trainer_id, group_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, categories))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(categories))
		}
	}

	return resultFunc
}

func GetExerciseInCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		//nullExercise := [0]ExerciseInCatetory{}
		category_seq := c.Param("category_seq")
		Category_Seq, err := strconv.Atoi(category_seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, category_seq))
			return
		}

		exercise := ExerciseInCatetory{}
		exercises, err := exercise.selectExerciseInCategory(Category_Seq, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, exercises))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(exercises))
		}
	}

	return resultFunc
}

func PostCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category := Category{}
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, category))
			return
		}

		_, err := common.IsDuplicatedTitle("t_category", category.Title, db)
		switch err {
		case nil:
			_, err := category.addCategory(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.FailedResponse(err, category))
			} else {
				c.JSON(http.StatusCreated, common.SucceedResponse(category))
			}

		case common.DuplicatedTitle:
			c.JSON(http.StatusConflict, common.FailedResponse(err, category.Title))

		default:
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, category))
		}
	}

	return resultFunc
}

func DeleteCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category_seq := c.Param("category_seq")
		Category_Seq, err := strconv.Atoi(category_seq)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, category_seq))
			return
		}

		row, err := deleteCategory(db, Category_Seq)
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

func PatchCategory(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		category := Category{}
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, category))
			return
		}

		row, err := common.IsDuplicatedTitle("t_category", category.Title, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, row))
		} else {
			switch row {
			case 0:
				row, err := category.modifyCategory(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, common.FailedResponse(err, row))
				} else {
					c.JSON(http.StatusCreated, common.SucceedResponse(row))
				}

			case 1:
				c.JSON(http.StatusConflict, common.SucceedResponse(category))
			}

		}

	}

	return resultFunc
}

func CategoryRouter(router *gin.Engine, db *sql.DB) {
	category := router.Group("/api/category")
	// GET All Category.
	// curl http://127.0.0.1:8080/api/category/all?trainer_id=Park&group_name=dygym -X GET
	category.GET("/all", GetAllCategory(db))

	// GET All Exercises in Specific Category.
	// curl http://127.0.0.1:8080/api/category/exercise/2?trainer_id=Park&group_name=dygym -X GET
	category.GET("/exercise/:category_seq", GetExerciseInCategory(db))

	// Create Category
	// curl http://127.0.0.1:8080/api/category -X POST -d '{"title": "등", "desc": "등 근육의 전반적인 향상", "group_name": "dygym", "trainer_id": "Choi Trainer","created_user": "Park", "updated_user": "Park"}' -H "Content-Type: application/json"
	category.POST("/", PostCategory(db))

	// Delete Specific Category.
	// curl http://127.0.0.1:8080/api/category/delete/5 -X DELETE
	category.DELETE("/:category_seq", DeleteCategory(db))

	// Update Specific Category.
	// curl http://127.0.0.1:8080/api/category/patch/16 -X PATCH -d {"title": "변경된 카테고리", "desc": "Blah Blah.."}
	category.PATCH("/", PatchCategory(db))
}

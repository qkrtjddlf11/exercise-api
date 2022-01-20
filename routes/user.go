package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type User struct {
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	Id           string  `json:"id"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Password     string  `json:"password"`
	Email        string  `json:"email"`
	Use_Yn       string  `json:"use_yn"`
}

type UserList struct {
	Seq        int    `json:"seq"`
	Name       string `json:"name"`
	Group_Name string `json:"group_name"`
	Trainer_Id string `json:"trainer_id"`
	Email      string `json:"email"`
	Use_Yn     string `json:"use_yn"`
}

func (u UserList) selectAllUser(db *sql.DB, trainer_id, group_name string) ([]UserList, error) {
	var users []UserList
	nullUsers := [0]UserList{}
	rows, err := db.Query(
		`SELECT 
			seq, 
			name,
			trainer_id,
			group_name,
			email,
			use_yn FROM t_user
		WHERE trainer_id = ? AND group_name = ?`, trainer_id, group_name)
	if err != nil {
		return nullUsers[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var user UserList
		rows.Scan(
			&user.Seq,
			&user.Name,
			&user.Trainer_Id,
			&user.Group_Name,
			&user.Email,
			&user.Use_Yn)
		users = append(users, user)
	}
	defer rows.Close()

	if len(users) > 0 {
		return users, nil
	} else {
		return nullUsers[:], errors.Wrap(err, "Users less than 1")
	}
}

func (u User) addUser(db *sql.DB) (int, error) {
	var id int
	stmt, err := db.Prepare(
		`INSERT INTO
			t_user(name,
				id,
				group_name,
				trainer_id,
				created_user,
				updated_user,
				email,
				use_yn)
			VALUE(?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		id = 0
		return id, errors.Wrap(err, "Failed to prepare to Database")
	}

	result, err := stmt.Exec(
		u.Name, u.Id, u.Group_Name, u.Trainer_Id, u.Trainer_Id, u.Trainer_Id, u.Email, u.Use_Yn)
	if err != nil {
		id = 0
		return id, errors.Wrap(err, "Failed to execute to Database")
	}

	seq, err := result.LastInsertId()
	if err != nil {
		id = 0
		return id, errors.Wrap(err, "Failed to insert last id to Database")
	}

	id = int(seq)
	defer stmt.Close()

	return id, nil
}

func GetAllUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		user := UserList{}
		users, err := user.selectAllUser(db, trainer_id, group_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, users))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(users))
		}
	}

	return resultFunc
}

func PostUser(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user := User{}
		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, user))
		}

		_, err := user.addUser(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, user))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(user))
		}
	}

	return resultFunc
}

func UserListRouter(router *gin.Engine, db *sql.DB) {
	user := router.Group("/api/user")

	user.GET("/all", GetAllUserList(db))
	user.POST("/", PostUser(db))
}

package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
type User struct {
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	Id           string  `json:"id"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User *string `json:"created_user"`
	Updated_User *string `json:"updated_user"`
	Use_Yn       *string `json:"use_yn"`
	Password     *string `json:"password"`
	Email        *string `json:"email"`
}
*/

type User struct {
	Seq          int    `json:"seq"`
	Name         string `json:"name"`
	Id           string `json:"id"`
	Group_Name   string `json:"group_name"`
	Trainer_Id   string `json:"trainer_id"`
	Created_Date string `json:"created_date"`
	Updated_Date string `json:"updated_date"`
	Created_User string `json:"created_user"`
	Updated_User string `json:"updated_user"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Use_Yn       string `json:"use_yn"`
}

func userOkResponse(data interface{}) gin.H {
	result := gin.H{
		"message": "",
		"status":  "ok",
		"result":  data,
	}

	return result
}

func userFailedResponse(err error, data interface{}) gin.H {
	result := gin.H{
		"message": err.Error(),
		"status":  "fail",
		"result":  data,
	}

	return result
}

func (u User) selectAllUser(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(
		`SELECT seq, 
		name,
		id,
		group_name,
		trainer_id,
		created_date,
		updated_date,
		created_user,
		updated_user,
		password,
		email,
		use_yn FROM t_user`)
	if err != nil {
		return
	}

	for rows.Next() {
		var user User
		rows.Scan(
			&u.Seq,
			&u.Name,
			&u.Id,
			&u.Password,
			&u.Email,
			&u.Group_Name,
			&u.Trainer_Id,
			&u.Created_Date,
			&u.Updated_Date,
			&u.Created_User,
			&u.Updated_User,
			&u.Use_Yn)
		users = append(users, user)
	}
	defer rows.Close()

	return
}

func getAllUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user := User{}
		users, err := user.selectAllUser(db)
		if err != nil {
			nullUsers := [0]User{}
			c.JSON(http.StatusInternalServerError, userFailedResponse(err, nullUsers))
		} else {
			c.JSON(http.StatusOK, userOkResponse(users))
		}
	}

	return resultFunc
}

func UserListRouter(router *gin.Engine, db *sql.DB) {
	user := router.Group("/api/user")

	user.GET("/all", getAllUserList(db))
}

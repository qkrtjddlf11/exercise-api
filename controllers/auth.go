package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type User struct {
	Seq          int     `json:"seq"`
	User_Name    string  `json:"user_name" binding:"required"`
	User_Id      string  `json:"user_id" binding:"required"`
	Group_Id     string  `json:"group_id" binding:"required"`
	Trainer_Id   string  `json:"trainer_id" binding:"required"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User *string `json:"created_user"`
	Updated_User *string `json:"updated_user"`
}

func (u User) insertUser(db *sql.DB) (Id int, err error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_user(user_name, user_id, group_id, trainer_id, created_user, updated_user) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	result, err := stmt.Exec(
		u.User_Name, u.User_Id, u.Group_Id, u.Trainer_Id, u.Created_User, u.Updated_User)
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

func postUser(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user := User{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"body": user,
				},
			})
			return
		}

		count, err := common.DuplicatedUserIdCheck("t_user", user.User_Id, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		} else {
			switch {
			case count == 0:
				_, err = user.insertUser(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, gin.H{})
				}
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("Duplicated user_id"),
					"parameters": gin.H{
						"body": user,
					},
				})
			}
		}

	}

	return resultFunc
}

func AuthRouter(router *gin.Engine, db *sql.DB) {
	public := router.Group("/api/auth")

	public.POST("/register", postUser(db))
}

package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

const (
	userId   = "parka"
	passWord = "1111"
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
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, user))
			return
		}

		count, err := common.DuplicatedUserIdCheck("t_user", user.User_Id, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, count))
		} else {
			switch {
			case count == 0:
				_, err = user.insertUser(db)
				if err != nil {
					c.JSON(http.StatusInternalServerError, common.FailedResponse(err, count))
				} else {
					c.JSON(http.StatusOK, common.SucceedResponse(user))
				}
			default:
				c.JSON(http.StatusBadRequest, common.FailedResponse(err, count))
			}
		}

	}

	return resultFunc
}

func loginUser(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		session := sessions.Default(c)
		user_id := c.PostForm("user_id")
		password := c.PostForm("password")
		if user_id != userId || password != passWord {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}

		session.Set(user_id, user_id)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	}

	return resultFunc
}

func logoutUser(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		session := sessions.Default(c)
		userSession := session.Get("park")
		if userSession == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
			return
		}

		session.Delete("park")
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}

	return resultFunc
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	userSesssion := session.Get("park")
	if userSesssion == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func AuthRouter(router *gin.Engine, db *sql.DB) {
	public := router.Group("/api/auth")

	public.POST("/register", postUser(db))

	store := cookie.NewStore([]byte("secretKey"))
	public.Use(sessions.Sessions("mysession", store))
	public.POST("/login", loginUser(db))
	public.GET("/logout", logoutUser(db))
	public.Use(authRequired)
	{
		public.GET("/me", me)
	}
}

func me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("parka")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/qkrtjddlf11/exercise-api/routes"
)

const MariaDB string = "mysql"

var db *sql.DB

func main() {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	//logFile, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(logFile)

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

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custrom log format
		logFormat := fmt.Sprintf("[API] %s - \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			//param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)

		logFile, _ := os.Create("gin.log")
		defer logFile.Close()

		log.SetOutput(logFile)
		log.Println(logFormat)

		return logFormat

	}))
	routes.CategoryRouter(router, db)
	routes.ExerciseRouter(router, db)
	routes.TodayExerciseRouter(router, db)

	router.Run("0.0.0.0:8080")
}

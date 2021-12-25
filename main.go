package main

import (
	"database/sql"
	"fmt"
	"io"
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

	routes.CategoryRouter(router, db)
	routes.ExerciseRouter(router, db)

	router.Run("0.0.0.0:8080")
}

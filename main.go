package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yasharya2901/smart_divide/database"
	"github.com/yasharya2901/smart_divide/models"
	"github.com/yasharya2901/smart_divide/routes"
)

func main() {
	// Connect to the database
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")

	db, err := database.NewMySQL(user, password, host, port, dbName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = db.GetDB().AutoMigrate(
		&models.Event{},
		&models.Expense{},
		&models.Person{},
		&models.ExpensePerson{},
	)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v0")

	routes.PersonRoutes(api)
	routes.EventRoutes(api, db.GetDB())
	routes.ExpenseRoutes(api)

	serverPort := os.Getenv("SERVER_PORT")

	if err := router.Run(fmt.Sprintf(":%s", serverPort)); err != nil {
		log.Fatal("Server Run Failed:", err)
	}

}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Migrate the schema
	err = db.GetDB().AutoMigrate(
		&models.Event{},
		&models.Expense{},
		&models.Person{},
		&models.ExpensePerson{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the server
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v0")

	routes.PersonRoutes(api, db.GetDB())
	routes.EventRoutes(api, db.GetDB())
	routes.ExpenseRoutes(api, db.GetDB())

	// Create http.Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: router,
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	// Start the server
	serverPort := os.Getenv("SERVER_PORT")

	log.Println("Starting server on port", serverPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

}

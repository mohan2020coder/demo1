package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

// CreateBook adds a new book
func (r *Repository) CreateBook(context *gin.Context) {
	var book Book
	if err := context.ShouldBindJSON(&book); err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Create(&book).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not create book"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "book has been added"})
}

// DeleteBook deletes a book by ID
func (r *Repository) DeleteBook(context *gin.Context) {
	id := context.Param("id")
	var book Book

	if err := r.DB.Where("id = ?", id).Delete(&book).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not delete book"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// GetBooks returns all books
func (r *Repository) GetBooks(context *gin.Context) {
	var books []Book

	if err := r.DB.Find(&books).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not get books"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "books fetched successfully", "data": books})
}

// GetBookByID returns a book by ID
func (r *Repository) GetBookByID(context *gin.Context) {
	id := context.Param("id")
	var book Book

	if err := r.DB.First(&book, id).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not get the book"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "book fetched successfully", "data": book})
}

// InsertDummyData inserts 5 dummy book records into the database
func (r *Repository) InsertDummyData() error {
	// Dummy data
	dummyBooks := []Book{
		{Author: "John Doe", Title: "Dummy Book 1", Publisher: "Publisher A"},
		{Author: "Jane Smith", Title: "Dummy Book 2", Publisher: "Publisher B"},
		{Author: "Alice Johnson", Title: "Dummy Book 3", Publisher: "Publisher C"},
		{Author: "Bob Brown", Title: "Dummy Book 4", Publisher: "Publisher D"},
		{Author: "Emma White", Title: "Dummy Book 5", Publisher: "Publisher E"},
	}

	// Insert dummy data into the database
	if err := r.DB.Create(&dummyBooks).Error; err != nil {
		return err
	}

	return nil
}

// SetupRoutes configures API routes
func (r *Repository) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/create_books", r.CreateBook)
		api.DELETE("/delete_book/:id", r.DeleteBook)
		api.GET("/get_books/:id", r.GetBookByID)
		api.GET("/books", r.GetBooks)
	}
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	// Establish database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	// Migrate the database schema
	if err := db.AutoMigrate(&Book{}); err != nil {
		log.Fatal("could not migrate database: ", err)
	}

	// Initialize repository
	repo := &Repository{DB: db}

	// Insert dummy data
	if err := repo.InsertDummyData(); err != nil {
		log.Fatal("could not insert dummy data: ", err)
	}

	// Initialize Gin router
	router := gin.Default()
	repo.SetupRoutes(router)

	// Start HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("could not start server: ", err)
	}
}

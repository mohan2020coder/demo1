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

func (r *Repository) DeleteBook(context *gin.Context) {
	id := context.Param("id")
	var book Book

	if err := r.DB.Where("id = ?", id).Delete(&book).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not delete book"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

func (r *Repository) GetBooks(context *gin.Context) {
	var books []Book

	if err := r.DB.Find(&books).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not get books"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "books fetched successfully", "data": books})
}

func (r *Repository) GetBookByID(context *gin.Context) {
	id := context.Param("id")
	var book Book

	if err := r.DB.First(&book, id).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not get the book"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "book fetched successfully", "data": book})
}

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
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("could not migrate database: ", err)
	}

	repo := &Repository{DB: db}

	router := gin.Default()
	repo.SetupRoutes(router)

	router.Run(":8080")
}

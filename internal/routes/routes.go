package routes

import (
	"library-management-go/internal/handlers"
	"library-management-go/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize services
	bookService := services.NewBookService(db)
	authorService := services.NewAuthorService(db)
	borrowerService := services.NewBorrowerService(db)
	borrowingService := services.NewBorrowingService(db)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService)
	authorHandler := handlers.NewAuthorHandler(authorService)
	borrowerHandler := handlers.NewBorrowerHandler(borrowerService)
	borrowingHandler := handlers.NewBorrowingHandler(borrowingService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Library Management API is running"})
		})

		// Author routes
		authors := v1.Group("/authors")
		{
			authors.POST("", authorHandler.CreateAuthor)
			authors.GET("", authorHandler.GetAllAuthors)
			authors.GET("/:id", authorHandler.GetAuthor)
			authors.PUT("/:id", authorHandler.UpdateAuthor)
			authors.DELETE("/:id", authorHandler.DeleteAuthor)
		}

		// Book routes
		books := v1.Group("/books")
		{
			books.POST("", bookHandler.CreateBook)
			books.GET("", bookHandler.GetAllBooks)
			books.GET("/:id", bookHandler.GetBook)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
		}

		// Borrower routes
		borrowers := v1.Group("/borrowers")
		{
			borrowers.POST("", borrowerHandler.CreateBorrower)
			borrowers.GET("", borrowerHandler.GetAllBorrowers)
			borrowers.GET("/:id", borrowerHandler.GetBorrower)
			borrowers.PUT("/:id", borrowerHandler.UpdateBorrower)
			borrowers.DELETE("/:id", borrowerHandler.DeleteBorrower)
		}

		// Borrowing routes
		borrowings := v1.Group("/borrowings")
		{
			borrowings.POST("/borrow", borrowingHandler.BorrowBook)
			borrowings.POST("/return", borrowingHandler.ReturnBook)
			borrowings.GET("", borrowingHandler.GetAllBorrowings)
			borrowings.GET("/:id", borrowingHandler.GetBorrowing)
			borrowings.GET("/borrower/:borrowerId", borrowingHandler.GetBorrowingsByBorrower)
			borrowings.GET("/overdue", borrowingHandler.GetOverdueBorrowings)
			borrowings.PUT("/update-overdue", borrowingHandler.UpdateOverdueStatus)
		}
	}
}

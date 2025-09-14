package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Author represents a book author
type Author struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Biography string    `json:"biography"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Book represents a book in the library
type Book struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	ISBN        string    `json:"isbn" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	AuthorID    uuid.UUID `json:"author_id" gorm:"type:uuid;not null"`
	Author      Author    `json:"author" gorm:"foreignKey:AuthorID"`
	PublishedAt time.Time `json:"published_at"`
	Available   bool      `json:"available" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Borrower represents a library member
type Borrower struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Borrowing represents a book borrowing record
type Borrowing struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookID     uuid.UUID `json:"book_id" gorm:"type:uuid;not null"`
	Book       Book      `json:"book" gorm:"foreignKey:BookID"`
	BorrowerID uuid.UUID `json:"borrower_id" gorm:"type:uuid;not null"`
	Borrower   Borrower  `json:"borrower" gorm:"foreignKey:BorrowerID"`
	BorrowedAt time.Time `json:"borrowed_at" gorm:"not null"`
	DueDate    time.Time `json:"due_date" gorm:"not null"`
	ReturnedAt *time.Time `json:"returned_at"`
	Status     string    `json:"status" gorm:"default:'borrowed'"` // borrowed, returned, overdue
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request DTOs
type CreateAuthorRequest struct {
	Name      string `json:"name" binding:"required"`
	Biography string `json:"biography"`
}

type UpdateAuthorRequest struct {
	Name      string `json:"name"`
	Biography string `json:"biography"`
}

type CreateBookRequest struct {
	Title       string    `json:"title" binding:"required"`
	ISBN        string    `json:"isbn" binding:"required"`
	Description string    `json:"description"`
	AuthorID    uuid.UUID `json:"author_id" binding:"required"`
	PublishedAt time.Time `json:"published_at"`
}

type UpdateBookRequest struct {
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	AuthorID    uuid.UUID `json:"author_id"`
	PublishedAt time.Time `json:"published_at"`
}

type CreateBorrowerRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdateBorrowerRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type BorrowBookRequest struct {
	BookID     uuid.UUID `json:"book_id" binding:"required"`
	BorrowerID uuid.UUID `json:"borrower_id" binding:"required"`
	DueDate    time.Time `json:"due_date" binding:"required"`
}

type ReturnBookRequest struct {
	BorrowingID uuid.UUID `json:"borrowing_id" binding:"required"`
}

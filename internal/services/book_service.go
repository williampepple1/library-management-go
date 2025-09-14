package services

import (
	"errors"

	"library-management-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.Book, error) {
	// Check if author exists
	var author models.Author
	if err := s.db.First(&author, req.AuthorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	// Check if ISBN already exists
	var existingBook models.Book
	if err := s.db.Where("isbn = ?", req.ISBN).First(&existingBook).Error; err == nil {
		return nil, errors.New("book with this ISBN already exists")
	}

	book := &models.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		Description: req.Description,
		AuthorID:    req.AuthorID,
		PublishedAt: req.PublishedAt,
		Available:   true,
	}

	if err := s.db.Create(book).Error; err != nil {
		return nil, err
	}

	// Load the author relationship
	if err := s.db.Preload("Author").First(book, book.ID).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetBook(id uuid.UUID) (*models.Book, error) {
	var book models.Book
	if err := s.db.Preload("Author").First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

func (s *BookService) GetAllBooks(page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&models.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get books with pagination
	if err := s.db.Preload("Author").
		Offset(offset).
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *BookService) UpdateBook(id uuid.UUID, req *models.UpdateBookRequest) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	// Check if author exists (if provided)
	if req.AuthorID != uuid.Nil {
		var author models.Author
		if err := s.db.First(&author, req.AuthorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("author not found")
			}
			return nil, err
		}
		book.AuthorID = req.AuthorID
	}

	// Check if ISBN already exists (if provided and different)
	if req.ISBN != "" && req.ISBN != book.ISBN {
		var existingBook models.Book
		if err := s.db.Where("isbn = ? AND id != ?", req.ISBN, id).First(&existingBook).Error; err == nil {
			return nil, errors.New("book with this ISBN already exists")
		}
		book.ISBN = req.ISBN
	}

	// Update fields
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if !req.PublishedAt.IsZero() {
		book.PublishedAt = req.PublishedAt
	}

	if err := s.db.Save(&book).Error; err != nil {
		return nil, err
	}

	// Load the author relationship
	if err := s.db.Preload("Author").First(&book, book.ID).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookService) DeleteBook(id uuid.UUID) error {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("book not found")
		}
		return err
	}

	// Check if book is currently borrowed
	var borrowing models.Borrowing
	if err := s.db.Where("book_id = ? AND status = 'borrowed'", id).First(&borrowing).Error; err == nil {
		return errors.New("cannot delete book that is currently borrowed")
	}

	if err := s.db.Delete(&book).Error; err != nil {
		return err
	}

	return nil
}

func (s *BookService) SearchBooks(query string, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Book{}).
		Joins("JOIN authors ON books.author_id = authors.id").
		Where("books.title ILIKE ? OR books.isbn ILIKE ? OR authors.name ILIKE ?",
			searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get books with pagination
	if err := s.db.Preload("Author").
		Joins("JOIN authors ON books.author_id = authors.id").
		Where("books.title ILIKE ? OR books.isbn ILIKE ? OR authors.name ILIKE ?",
			searchQuery, searchQuery, searchQuery).
		Offset(offset).
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

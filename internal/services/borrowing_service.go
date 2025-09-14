package services

import (
	"errors"
	"time"

	"library-management-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BorrowingService struct {
	db *gorm.DB
}

func NewBorrowingService(db *gorm.DB) *BorrowingService {
	return &BorrowingService{db: db}
}

func (s *BorrowingService) BorrowBook(req *models.BorrowBookRequest) (*models.Borrowing, error) {
	// Check if book exists and is available
	var book models.Book
	if err := s.db.First(&book, req.BookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if !book.Available {
		return nil, errors.New("book is not available for borrowing")
	}

	// Check if borrower exists
	var borrower models.Borrower
	if err := s.db.First(&borrower, req.BorrowerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrower not found")
		}
		return nil, err
	}

	// Check if borrower has any overdue books
	var overdueCount int64
	if err := s.db.Model(&models.Borrowing{}).
		Where("borrower_id = ? AND status = 'borrowed' AND due_date < ?", req.BorrowerID, time.Now()).
		Count(&overdueCount).Error; err != nil {
		return nil, err
	}

	if overdueCount > 0 {
		return nil, errors.New("borrower has overdue books and cannot borrow new books")
	}

	// Check if borrower has reached maximum borrowing limit (e.g., 5 books)
	var activeBorrowingCount int64
	if err := s.db.Model(&models.Borrowing{}).
		Where("borrower_id = ? AND status = 'borrowed'", req.BorrowerID).
		Count(&activeBorrowingCount).Error; err != nil {
		return nil, err
	}

	if activeBorrowingCount >= 5 {
		return nil, errors.New("borrower has reached maximum borrowing limit")
	}

	// Create borrowing record
	borrowing := &models.Borrowing{
		BookID:     req.BookID,
		BorrowerID: req.BorrowerID,
		BorrowedAt: time.Now(),
		DueDate:    req.DueDate,
		Status:     "borrowed",
	}

	if err := s.db.Create(borrowing).Error; err != nil {
		return nil, err
	}

	// Update book availability
	book.Available = false
	if err := s.db.Save(&book).Error; err != nil {
		return nil, err
	}

	// Load relationships
	if err := s.db.Preload("Book.Author").Preload("Borrower").First(borrowing, borrowing.ID).Error; err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (s *BorrowingService) ReturnBook(req *models.ReturnBookRequest) (*models.Borrowing, error) {
	var borrowing models.Borrowing
	if err := s.db.Preload("Book").Preload("Borrower").First(&borrowing, req.BorrowingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrowing record not found")
		}
		return nil, err
	}

	if borrowing.Status != "borrowed" {
		return nil, errors.New("book is not currently borrowed")
	}

	// Update borrowing record
	now := time.Now()
	borrowing.ReturnedAt = &now
	borrowing.Status = "returned"

	// Check if overdue
	if now.After(borrowing.DueDate) {
		borrowing.Status = "overdue"
	}

	if err := s.db.Save(&borrowing).Error; err != nil {
		return nil, err
	}

	// Update book availability
	borrowing.Book.Available = true
	if err := s.db.Save(&borrowing.Book).Error; err != nil {
		return nil, err
	}

	return &borrowing, nil
}

func (s *BorrowingService) GetBorrowing(id uuid.UUID) (*models.Borrowing, error) {
	var borrowing models.Borrowing
	if err := s.db.Preload("Book.Author").Preload("Borrower").First(&borrowing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrowing record not found")
		}
		return nil, err
	}
	return &borrowing, nil
}

func (s *BorrowingService) GetAllBorrowings(page, limit int) ([]models.Borrowing, int64, error) {
	var borrowings []models.Borrowing
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&models.Borrowing{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get borrowings with pagination
	if err := s.db.Preload("Book.Author").Preload("Borrower").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&borrowings).Error; err != nil {
		return nil, 0, err
	}

	return borrowings, total, nil
}

func (s *BorrowingService) GetBorrowingsByBorrower(borrowerID uuid.UUID, page, limit int) ([]models.Borrowing, int64, error) {
	var borrowings []models.Borrowing
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&models.Borrowing{}).Where("borrower_id = ?", borrowerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get borrowings with pagination
	if err := s.db.Preload("Book.Author").Preload("Borrower").
		Where("borrower_id = ?", borrowerID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&borrowings).Error; err != nil {
		return nil, 0, err
	}

	return borrowings, total, nil
}

func (s *BorrowingService) GetOverdueBorrowings(page, limit int) ([]models.Borrowing, int64, error) {
	var borrowings []models.Borrowing
	var total int64

	offset := (page - 1) * limit
	now := time.Now()

	// Count total records
	if err := s.db.Model(&models.Borrowing{}).
		Where("status = 'borrowed' AND due_date < ?", now).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get overdue borrowings with pagination
	if err := s.db.Preload("Book.Author").Preload("Borrower").
		Where("status = 'borrowed' AND due_date < ?", now).
		Offset(offset).Limit(limit).
		Order("due_date ASC").
		Find(&borrowings).Error; err != nil {
		return nil, 0, err
	}

	return borrowings, total, nil
}

func (s *BorrowingService) UpdateOverdueStatus() error {
	now := time.Now()
	
	// Update borrowings that are now overdue
	if err := s.db.Model(&models.Borrowing{}).
		Where("status = 'borrowed' AND due_date < ?", now).
		Update("status", "overdue").Error; err != nil {
		return err
	}

	return nil
}

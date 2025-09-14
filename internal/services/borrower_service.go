package services

import (
	"errors"

	"library-management-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BorrowerService struct {
	db *gorm.DB
}

func NewBorrowerService(db *gorm.DB) *BorrowerService {
	return &BorrowerService{db: db}
}

func (s *BorrowerService) CreateBorrower(req *models.CreateBorrowerRequest) (*models.Borrower, error) {
	// Check if email already exists
	var existingBorrower models.Borrower
	if err := s.db.Where("email = ?", req.Email).First(&existingBorrower).Error; err == nil {
		return nil, errors.New("borrower with this email already exists")
	}

	borrower := &models.Borrower{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err := s.db.Create(borrower).Error; err != nil {
		return nil, err
	}

	return borrower, nil
}

func (s *BorrowerService) GetBorrower(id uuid.UUID) (*models.Borrower, error) {
	var borrower models.Borrower
	if err := s.db.First(&borrower, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrower not found")
		}
		return nil, err
	}
	return &borrower, nil
}

func (s *BorrowerService) GetAllBorrowers(page, limit int) ([]models.Borrower, int64, error) {
	var borrowers []models.Borrower
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&models.Borrower{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get borrowers with pagination
	if err := s.db.Offset(offset).Limit(limit).Find(&borrowers).Error; err != nil {
		return nil, 0, err
	}

	return borrowers, total, nil
}

func (s *BorrowerService) UpdateBorrower(id uuid.UUID, req *models.UpdateBorrowerRequest) (*models.Borrower, error) {
	var borrower models.Borrower
	if err := s.db.First(&borrower, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrower not found")
		}
		return nil, err
	}

	// Check if email already exists (if provided and different)
	if req.Email != "" && req.Email != borrower.Email {
		var existingBorrower models.Borrower
		if err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existingBorrower).Error; err == nil {
			return nil, errors.New("borrower with this email already exists")
		}
		borrower.Email = req.Email
	}

	// Update fields
	if req.Name != "" {
		borrower.Name = req.Name
	}
	if req.Phone != "" {
		borrower.Phone = req.Phone
	}
	if req.Address != "" {
		borrower.Address = req.Address
	}

	if err := s.db.Save(&borrower).Error; err != nil {
		return nil, err
	}

	return &borrower, nil
}

func (s *BorrowerService) DeleteBorrower(id uuid.UUID) error {
	var borrower models.Borrower
	if err := s.db.First(&borrower, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("borrower not found")
		}
		return err
	}

	// Check if borrower has active borrowings
	var borrowingCount int64
	if err := s.db.Model(&models.Borrowing{}).Where("borrower_id = ? AND status = 'borrowed'", id).Count(&borrowingCount).Error; err != nil {
		return err
	}

	if borrowingCount > 0 {
		return errors.New("cannot delete borrower with active borrowings")
	}

	if err := s.db.Delete(&borrower).Error; err != nil {
		return err
	}

	return nil
}

func (s *BorrowerService) SearchBorrowers(query string, page, limit int) ([]models.Borrower, int64, error) {
	var borrowers []models.Borrower
	var total int64

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Borrower{}).
		Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get borrowers with pagination
	if err := s.db.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchQuery, searchQuery, searchQuery).
		Offset(offset).
		Limit(limit).
		Find(&borrowers).Error; err != nil {
		return nil, 0, err
	}

	return borrowers, total, nil
}

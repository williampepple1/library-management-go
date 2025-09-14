package services

import (
	"errors"

	"library-management-go/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorService struct {
	db *gorm.DB
}

func NewAuthorService(db *gorm.DB) *AuthorService {
	return &AuthorService{db: db}
}

func (s *AuthorService) CreateAuthor(req *models.CreateAuthorRequest) (*models.Author, error) {
	author := &models.Author{
		Name:      req.Name,
		Biography: req.Biography,
	}

	if err := s.db.Create(author).Error; err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) GetAuthor(id uuid.UUID) (*models.Author, error) {
	var author models.Author
	if err := s.db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}
	return &author, nil
}

func (s *AuthorService) GetAllAuthors(page, limit int) ([]models.Author, int64, error) {
	var authors []models.Author
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&models.Author{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get authors with pagination
	if err := s.db.Offset(offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

func (s *AuthorService) UpdateAuthor(id uuid.UUID, req *models.UpdateAuthorRequest) (*models.Author, error) {
	var author models.Author
	if err := s.db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		author.Name = req.Name
	}
	if req.Biography != "" {
		author.Biography = req.Biography
	}

	if err := s.db.Save(&author).Error; err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *AuthorService) DeleteAuthor(id uuid.UUID) error {
	var author models.Author
	if err := s.db.First(&author, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("author not found")
		}
		return err
	}

	// Check if author has books
	var bookCount int64
	if err := s.db.Model(&models.Book{}).Where("author_id = ?", id).Count(&bookCount).Error; err != nil {
		return err
	}

	if bookCount > 0 {
		return errors.New("cannot delete author with existing books")
	}

	if err := s.db.Delete(&author).Error; err != nil {
		return err
	}

	return nil
}

func (s *AuthorService) SearchAuthors(query string, page, limit int) ([]models.Author, int64, error) {
	var authors []models.Author
	var total int64

	offset := (page - 1) * limit
	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Author{}).
		Where("name ILIKE ? OR biography ILIKE ?", searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get authors with pagination
	if err := s.db.Where("name ILIKE ? OR biography ILIKE ?", searchQuery, searchQuery).
		Offset(offset).
		Limit(limit).
		Find(&authors).Error; err != nil {
		return nil, 0, err
	}

	return authors, total, nil
}

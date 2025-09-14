package handlers

import (
	"net/http"
	"strconv"

	"library-management-go/internal/models"
	"library-management-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BorrowerHandler struct {
	borrowerService *services.BorrowerService
}

func NewBorrowerHandler(borrowerService *services.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{borrowerService: borrowerService}
}

func (h *BorrowerHandler) CreateBorrower(c *gin.Context) {
	var req models.CreateBorrowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrower, err := h.borrowerService.CreateBorrower(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": borrower})
}

func (h *BorrowerHandler) GetBorrower(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower ID"})
		return
	}

	borrower, err := h.borrowerService.GetBorrower(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": borrower})
}

func (h *BorrowerHandler) GetAllBorrowers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var borrowers []models.Borrower
	var total int64
	var err error

	if search != "" {
		borrowers, total, err = h.borrowerService.SearchBorrowers(search, page, limit)
	} else {
		borrowers, total, err = h.borrowerService.GetAllBorrowers(page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": borrowers,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *BorrowerHandler) UpdateBorrower(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower ID"})
		return
	}

	var req models.UpdateBorrowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrower, err := h.borrowerService.UpdateBorrower(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": borrower})
}

func (h *BorrowerHandler) DeleteBorrower(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower ID"})
		return
	}

	err = h.borrowerService.DeleteBorrower(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "borrower deleted successfully"})
}

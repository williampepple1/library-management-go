package handlers

import (
	"net/http"
	"strconv"

	"library-management-go/internal/models"
	"library-management-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BorrowingHandler struct {
	borrowingService *services.BorrowingService
}

func NewBorrowingHandler(borrowingService *services.BorrowingService) *BorrowingHandler {
	return &BorrowingHandler{borrowingService: borrowingService}
}

func (h *BorrowingHandler) BorrowBook(c *gin.Context) {
	var req models.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrowing, err := h.borrowingService.BorrowBook(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": borrowing})
}

func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
	var req models.ReturnBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrowing, err := h.borrowingService.ReturnBook(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": borrowing})
}

func (h *BorrowingHandler) GetBorrowing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrowing ID"})
		return
	}

	borrowing, err := h.borrowingService.GetBorrowing(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": borrowing})
}

func (h *BorrowingHandler) GetAllBorrowings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	borrowings, total, err := h.borrowingService.GetAllBorrowings(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": borrowings,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *BorrowingHandler) GetBorrowingsByBorrower(c *gin.Context) {
	borrowerIDStr := c.Param("borrowerId")
	borrowerID, err := uuid.Parse(borrowerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	borrowings, total, err := h.borrowingService.GetBorrowingsByBorrower(borrowerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": borrowings,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *BorrowingHandler) GetOverdueBorrowings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	borrowings, total, err := h.borrowingService.GetOverdueBorrowings(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": borrowings,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *BorrowingHandler) UpdateOverdueStatus(c *gin.Context) {
	err := h.borrowingService.UpdateOverdueStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "overdue status updated successfully"})
}

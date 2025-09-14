package handlers

import (
	"net/http"
	"strconv"

	"library-management-go/internal/models"
	"library-management-go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorHandler struct {
	authorService *services.AuthorService
}

func NewAuthorHandler(authorService *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{authorService: authorService}
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var req models.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.authorService.CreateAuthor(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": author})
}

func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author ID"})
		return
	}

	author, err := h.authorService.GetAuthor(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": author})
}

func (h *AuthorHandler) GetAllAuthors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var authors []models.Author
	var total int64
	var err error

	if search != "" {
		authors, total, err = h.authorService.SearchAuthors(search, page, limit)
	} else {
		authors, total, err = h.authorService.GetAllAuthors(page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": authors,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author ID"})
		return
	}

	var req models.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	author, err := h.authorService.UpdateAuthor(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": author})
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author ID"})
		return
	}

	err = h.authorService.DeleteAuthor(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "author deleted successfully"})
}

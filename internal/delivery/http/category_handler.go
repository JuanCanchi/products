package http

import (
	"github.com/gin-gonic/gin"
	"github.com/juancanchi/products/internal/domain"
	"github.com/juancanchi/products/internal/usecase"
	"net/http"
)

type CategoryHandler struct {
	usecase *usecase.CategoryUsecase
}

func NewCategoryHandler(u *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: u}
}

// POST /categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var cat domain.Category
	if err := c.ShouldBindJSON(&cat); err != nil || cat.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nombre requerido"})
		return
	}

	if err := h.usecase.Create(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear categoría"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// GET /categories
func (h *CategoryHandler) List(c *gin.Context) {
	cats, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener categorías"})
		return
	}

	c.JSON(http.StatusOK, cats)
}

// PUT /categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var cat domain.Category
	if err := c.ShouldBindJSON(&cat); err != nil || cat.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nombre requerido"})
		return
	}

	cat.ID = id
	if err := h.usecase.Update(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DELETE /categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar"})
		return
	}

	c.Status(http.StatusNoContent)
}

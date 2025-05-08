package http

import (
	"github.com/juancanchi/products/internal/domain"
	"github.com/juancanchi/products/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(u *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id no encontrado en el token"})
		return
	}
	p.UserID = userID.(string)
	p.Status = "PENDING"

	if err := h.usecase.Create(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) ListByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id no encontrado"})
		return
	}

	products, err := h.usecase.ListByUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p.ID = id // aseguramos que se use el ID de la URL

	if p.CategoryID != nil && *p.CategoryID == "" {
		p.CategoryID = nil
	}

	if err := h.usecase.Update(c.Request.Context(), &p, userID.(string)); err != nil {
		if err.Error() == "no autorizado" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no sos el dueño del producto"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	err := h.usecase.Delete(c.Request.Context(), id, userID.(string))
	if err != nil {
		if err.Error() == "no autorizado" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no sos el dueño del producto"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ChangeStatus(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "no autorizado"})
		return
	}

	id := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "estado inválido"})
		return
	}

	err := h.usecase.UpdateStatus(c.Request.Context(), id, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el estado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleListProducts godoc
// @Summary      List all products
// @Description  Returns all product names. Requires admin privileges.
// @Tags         products
// @Produce      json
// @Success      200  {array}   string
// @Failure      401  {string}  string             "Admin privileges required"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /product/list [get]
func (h *ProductHandler) HandleListProducts(c *gin.Context) {
	admin, exists := c.Get("admin")
	if !exists || !admin.(bool) {
		c.String(http.StatusUnauthorized, "Only admin allowed to list products")
		return
	}

	names, err := h.Repo.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, names)
}

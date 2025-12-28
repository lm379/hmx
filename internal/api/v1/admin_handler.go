package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleDeleteOpera (DELETE /api/v1/admin/operas/:id)
func HandleDeleteOpera(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.DeleteOpera(uint(operaID))
	if status >= 400 {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

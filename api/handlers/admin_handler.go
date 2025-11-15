package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
)

// HandleDeleteOpera (DELETE /api/v1/admin/operas/:id)
func HandleDeleteOpera(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid opera ID")
		return
	}

	response, status := services.DeleteOpera(uint(operaID))
	if status >= 400 {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}

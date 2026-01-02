package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/converter"
	"github.com/lm379/hmx/pkg/pagination"
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

// HandleAdminGetOperas (GET /api/v1/admin/operas)
func HandleAdminGetOperas(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	// 获取所有视频，包括隐藏的
	operas, total, err := services.GetOperas(pagination, true)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch operas")
		return
	}

	operaResponses := converter.ToOperaResponseList(operas)
	
	// Count stats
	if len(operaResponses) > 0 {
		operaIDs := make([]uint, len(operaResponses))
		for i, r := range operaResponses {
			operaIDs[i] = r.OperaID
		}

		likes, favorites, shares, plays := services.BatchGetCounts(operaIDs)
		for _, r := range operaResponses {
			r.LikeCount = likes[r.OperaID]
			r.FavoriteCount = favorites[r.OperaID]
			r.ShareCount = shares[r.OperaID]
			r.PlayCount = plays[r.OperaID]
		}
	}

	resp.Success(c, gin.H{
		"list": operaResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleAdminUpdateOpera (PUT /api/v1/admin/operas/:id)
func HandleAdminUpdateOpera(c *gin.Context) {
	idParam := c.Param("id")
	operaID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	var input models.UpdateOperaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.UpdateOpera(uint(operaID), input)
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminCreateArtist (POST /api/v1/admin/artists)
func HandleAdminCreateArtist(c *gin.Context) {
	var input models.CreateArtistRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.CreateArtist(input)
	if status != http.StatusCreated {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Created(c, response)
}

// HandleAdminUpdateArtist (PUT /api/v1/admin/artists/:id)
func HandleAdminUpdateArtist(c *gin.Context) {
	idParam := c.Param("id")
	artistID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid artist ID")
		return
	}

	var input models.UpdateArtistRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	response, status := services.UpdateArtist(uint(artistID), input)
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminDeleteArtist (DELETE /api/v1/admin/artists/:id)
func HandleAdminDeleteArtist(c *gin.Context) {
	idParam := c.Param("id")
	artistID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid artist ID")
		return
	}

	response, status := services.DeleteArtist(uint(artistID))
	if status != http.StatusOK {
		resp.Error(c, status, response["error"].(string))
		return
	}
	resp.Success(c, response)
}

// HandleAdminGetUsers (GET /api/v1/admin/users)
func HandleAdminGetUsers(c *gin.Context) {
	pagination := pagination.GetPagination(c)

	users, total, err := services.GetAllUsers(pagination)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch users")
		return
	}

	userResponses := make([]*models.UserResponse, len(users))
	for i, u := range users {
		// Create a copy of the loop variable to avoid pointing to the same address
		user := u
		userResponses[i] = converter.ToUserResponse(&user)
	}

	resp.Success(c, gin.H{
		"list": userResponses,
		"pagination": gin.H{
			"total":     total,
			"page":      pagination.Page,
			"page_size": pagination.PageSize,
		},
	})
}

// HandleAdminUpdateUserRole (PUT /api/v1/admin/users/:id/role)
func HandleAdminUpdateUserRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid user ID")
		return
	}

	var input models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdateUserRole(uint(userID), input.Role); err != nil {
		resp.InternalServerError(c, "Failed to update role")
		return
	}
	resp.Success(c, gin.H{"message": "Role updated successfully"})
}

// HandleAdminGetDashboardStats (GET /api/v1/admin/stats)
func HandleAdminGetDashboardStats(c *gin.Context) {
	stats, err := services.GetDashboardStats()
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch dashboard stats")
		return
	}
	resp.Success(c, stats)
}

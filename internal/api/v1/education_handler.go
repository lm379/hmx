package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/services"
	"github.com/lm379/hmx/pkg/pagination"
	resp "github.com/lm379/hmx/pkg/response"
	"gorm.io/gorm"
)

// HandleGetEducationBooks 返回前台可见的黄梅教育 PDF 列表。
func HandleGetEducationBooks(c *gin.Context) {
	p := pagination.GetPagination(c)
	educationService := services.NewEducationService()
	list, total, err := educationService.GetPublishedBooks(p)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch education books")
		return
	}
	resp.Success(c, gin.H{
		"list": list,
		"pagination": gin.H{
			"total":     total,
			"page":      p.Page,
			"page_size": p.PageSize,
		},
	})
}

// HandleGetEducationBookByID 返回单本已发布 PDF 的阅读地址。
func HandleGetEducationBookByID(c *gin.Context) {
	bookID, ok := parseBookID(c)
	if !ok {
		return
	}
	educationService := services.NewEducationService()
	detail, err := educationService.GetPublishedBookDetail(bookID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Education book not found")
			return
		}
		resp.InternalServerError(c, "Failed to fetch education book")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminGetEducationBooks 返回后台教育资源列表，包含草稿。
func HandleAdminGetEducationBooks(c *gin.Context) {
	p := pagination.GetPagination(c)
	educationService := services.NewEducationService()
	list, total, err := educationService.GetAllBooks(p)
	if err != nil {
		resp.InternalServerError(c, "Failed to fetch education books")
		return
	}
	resp.Success(c, gin.H{
		"list": list,
		"pagination": gin.H{
			"total":     total,
			"page":      p.Page,
			"page_size": p.PageSize,
		},
	})
}

// HandleAdminGetEducationBookByID 返回后台编辑所需的 PDF 资源详情。
func HandleAdminGetEducationBookByID(c *gin.Context) {
	bookID, ok := parseBookID(c)
	if !ok {
		return
	}
	educationService := services.NewEducationService()
	detail, err := educationService.GetAdminBookDetail(bookID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Education book not found")
			return
		}
		resp.InternalServerError(c, "Failed to fetch education book")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminCreateEducationBook 创建教育 PDF 资源。
func HandleAdminCreateEducationBook(c *gin.Context) {
	var input models.CreateEducationBookRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}
	educationService := services.NewEducationService()
	detail, err := educationService.CreateBook(input)
	if err != nil {
		resp.InternalServerError(c, "Failed to create education book")
		return
	}
	resp.Created(c, detail)
}

// HandleAdminUpdateEducationBook 更新教育 PDF 资源。
func HandleAdminUpdateEducationBook(c *gin.Context) {
	bookID, ok := parseBookID(c)
	if !ok {
		return
	}
	var input models.UpdateEducationBookRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		resp.BadRequest(c, err.Error())
		return
	}
	educationService := services.NewEducationService()
	detail, err := educationService.UpdateBook(bookID, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.NotFound(c, "Education book not found")
			return
		}
		resp.InternalServerError(c, "Failed to update education book")
		return
	}
	resp.Success(c, detail)
}

// HandleAdminDeleteEducationBook 删除教育 PDF 资源。
func HandleAdminDeleteEducationBook(c *gin.Context) {
	bookID, ok := parseBookID(c)
	if !ok {
		return
	}
	educationService := services.NewEducationService()
	if err := educationService.DeleteBook(bookID); err != nil {
		resp.InternalServerError(c, "Failed to delete education book")
		return
	}
	resp.Success(c, gin.H{"message": "Education book deleted"})
}

// parseBookID 统一解析路由中的教育资源 ID，并直接写出 400 响应。
func parseBookID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid education book ID")
		return 0, false
	}
	return uint(id), true
}

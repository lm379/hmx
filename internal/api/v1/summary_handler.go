package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/internal/services"
	resp "github.com/lm379/hmx/pkg/response"
)

// HandleGetVideoSummary 获取视频AI字幕摘要
// GET /api/v1/operas/:id/summary
func HandleGetVideoSummary(c *gin.Context) {
	operaIDStr := c.Param("id")
	operaID, err := strconv.ParseUint(operaIDStr, 10, 32)
	if err != nil {
		resp.BadRequest(c, "Invalid opera ID")
		return
	}

	summary, err := services.GetVideoSubtitleSummary(uint(operaID))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, gin.H{
		"opera_id": operaID,
		"summary":  summary,
		"status":   getSummaryStatus(summary),
	})
}

func getSummaryStatus(summary string) string {
	if summary == "正在生成中" {
		return "generating"
	}
	return "completed"
}

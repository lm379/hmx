package services

import (
	"log"
	"time"

	"github.com/lm379/hmx/database"
	"github.com/lm379/hmx/internal/models"
)

// TriggerAISummary 异步触发 AI 总结 (模拟)
// 在真实项目中，这里可能会调用一个外部 AI API
func TriggerAISummary(operaID uint) {
	log.Printf("[AI Service] 任务开始: 正在为 Opera ID %d 生成摘要...", operaID)

	// 模拟 AI 处理所需的时间
	time.Sleep(15 * time.Second) // 模拟 15 秒的 AI 处理

	// 模拟的 AI 结果
	mockSummary := "这是一段由 AI 模拟生成的黄梅戏摘要。这段摘要详细分析了...（此处为模拟内容）"

	// 将结果写回数据库
	db := database.DB
	var opera models.Opera
	if err := db.First(&opera, operaID).Error; err != nil {
		log.Printf("[AI Service] 错误: 找不到 Opera ID %d: %v", operaID, err)
		return
	}

	if err := db.Model(&opera).Update("ai_summary", mockSummary).Error; err != nil {
		log.Printf("[AI Service] 错误: 更新 Opera ID %d 摘要失败: %v", operaID, err)
		return
	}

	log.Printf("[AI Service] 任务完成: Opera ID %d 摘要已更新。", operaID)
}

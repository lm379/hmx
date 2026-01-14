package search

import (
	"log"
	"sync"

	"github.com/lm379/hmx/config"
	"github.com/meilisearch/meilisearch-go"
)

var (
	client    meilisearch.ServiceManager
	once      sync.Once
	clientErr error
)

// GetMeilisearchClient 获取 Meilisearch 客户端单例
func GetMeilisearchClient() (meilisearch.ServiceManager, error) {
	once.Do(func() {
		if config.AppConfig.MeilisearchHost == "" {
			log.Println("Meilisearch host not configured, using default: http://localhost:7700")
		}

		client = meilisearch.New(
			config.AppConfig.MeilisearchHost,
			meilisearch.WithAPIKey(config.AppConfig.MeilisearchAPIKey),
		)

		// 测试连接
		_, err := client.Health()
		if err != nil {
			log.Printf("Warning: Failed to connect to Meilisearch: %v", err)
			clientErr = err
			return
		}

		log.Println("Successfully connected to Meilisearch")
	})

	return client, clientErr
}

// InitializeIndexes 初始化所有索引
func InitializeIndexes() error {
	client, err := GetMeilisearchClient()
	if err != nil {
		return err
	}

	// 初始化 Opera 索引
	if err := initOperaIndex(client); err != nil {
		return err
	}

	// 初始化 Artist 索引
	if err := initArtistIndex(client); err != nil {
		return err
	}

	log.Println("All search indexes initialized successfully")
	return nil
}

// stringsToInterfaces 将字符串切片转换为interface{}切片
func stringsToInterfaces(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// initOperaIndex 初始化 Opera 索引
func initOperaIndex(client meilisearch.ServiceManager) error {
	indexName := "operas"

	// 创建或获取索引
	_, err := client.GetIndex(indexName)
	if err != nil {
		// 索引不存在，创建它
		_, err = client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        indexName,
			PrimaryKey: "opera_id",
		})
		if err != nil {
			log.Printf("Failed to create index %s: %v", indexName, err)
			return err
		}
		log.Printf("Created index: %s", indexName)
	}

	index := client.Index(indexName)

	// 配置可搜索属性：包含中文、拼音和分词字段
	attrs := []string{
		"opera_title",
		"opera_title_tokens", // 分词搜索（优先用于模糊匹配）
		"opera_title_py",     // 拼音搜索
		"artist_names",
		"artist_names_tokens", // 艺术家分词搜索
		"artist_names_py",     // 艺术家拼音搜索
	}
	_, err = index.UpdateSearchableAttributes(&attrs)
	if err != nil {
		log.Printf("Failed to update searchable attributes: %v", err)
		return err
	}

	// 配置可过滤属性
	filterAttrs := stringsToInterfaces([]string{
		"opera_id",
		"artist_ids",
	})
	_, err = index.UpdateFilterableAttributes(&filterAttrs)
	if err != nil {
		log.Printf("Failed to update filterable attributes: %v", err)
		return err
	}

	// 配置可排序属性
	sortAttrs := []string{
		"opera_id",
	}
	_, err = index.UpdateSortableAttributes(&sortAttrs)
	if err != nil {
		log.Printf("Failed to update sortable attributes: %v", err)
		return err
	}

	// 配置拼写容错（模糊搜索）- 针对中文优化
	typoToleranceSettings := &meilisearch.TypoTolerance{
		Enabled: true,
		MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
			OneTypo:  1, // 中文通常是单字或双字，1个字符就允许1个拼写错误
			TwoTypos: 3, // 3个字符以上允许2个拼写错误
		},
		DisableOnWords:      []string{}, // 不禁用任何词
		DisableOnAttributes: []string{}, // 不禁用任何属性
	}
	_, err = index.UpdateTypoTolerance(typoToleranceSettings)
	if err != nil {
		log.Printf("Warning: Failed to update typo tolerance: %v", err)
		// 不返回错误，因为这不是致命问题
	}

	// 配置搜索策略 - 优先精确匹配，然后是模糊匹配
	rankingRules := []string{
		"words",     // 匹配的词数
		"typo",      // 拼写错误数（越少越好）
		"proximity", // 词之间的距离
		"attribute", // 字段权重
		"sort",      // 排序
		"exactness", // 精确度
	}
	_, err = index.UpdateRankingRules(&rankingRules)
	if err != nil {
		log.Printf("Warning: Failed to update ranking rules: %v", err)
	}

	log.Printf("Configured index: %s with fuzzy search enabled", indexName)
	return nil
}

// initArtistIndex 初始化 Artist 索引
func initArtistIndex(client meilisearch.ServiceManager) error {
	indexName := "artists"

	// 创建或获取索引
	_, err := client.GetIndex(indexName)
	if err != nil {
		// 索引不存在，创建它
		_, err = client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        indexName,
			PrimaryKey: "artist_id",
		})
		if err != nil {
			log.Printf("Failed to create index %s: %v", indexName, err)
			return err
		}
		log.Printf("Created index: %s", indexName)
	}

	index := client.Index(indexName)

	// 配置可搜索属性：包含中文、拼音和分词字段
	artistAttrs := []string{
		"name",
		"name_tokens", // 分词搜索（优先用于模糊匹配）
		"name_py",     // 拼音搜索
	}
	_, err = index.UpdateSearchableAttributes(&artistAttrs)
	if err != nil {
		log.Printf("Failed to update searchable attributes: %v", err)
		return err
	}

	// 配置可过滤属性
	artistFilterAttrs := stringsToInterfaces([]string{
		"artist_id",
	})
	_, err = index.UpdateFilterableAttributes(&artistFilterAttrs)
	if err != nil {
		log.Printf("Failed to update filterable attributes: %v", err)
		return err
	}

	// 配置可排序属性
	artistSortAttrs := []string{
		"artist_id",
		"name",
	}
	_, err = index.UpdateSortableAttributes(&artistSortAttrs)
	if err != nil {
		log.Printf("Failed to update sortable attributes: %v", err)
		return err
	}

	// 配置拼写容错（模糊搜索）- 针对中文优化
	typoToleranceSettings := &meilisearch.TypoTolerance{
		Enabled: true,
		MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
			OneTypo:  1, // 中文名字通常很短，1个字符就允许1个拼写错误
			TwoTypos: 2, // 2个字符以上允许2个拼写错误
		},
		DisableOnWords:      []string{},
		DisableOnAttributes: []string{},
	}
	_, err = index.UpdateTypoTolerance(typoToleranceSettings)
	if err != nil {
		log.Printf("Warning: Failed to update typo tolerance: %v", err)
	}

	// 配置搜索策略
	rankingRules := []string{
		"words",
		"typo",
		"proximity",
		"attribute",
		"sort",
		"exactness",
	}
	_, err = index.UpdateRankingRules(&rankingRules)
	if err != nil {
		log.Printf("Warning: Failed to update ranking rules: %v", err)
	}

	log.Printf("Configured index: %s with fuzzy search enabled", indexName)
	return nil
}

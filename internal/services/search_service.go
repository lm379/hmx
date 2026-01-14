package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lm379/hmx/internal/models"
	"github.com/lm379/hmx/internal/repository"
	"github.com/lm379/hmx/pkg/pinyin"
	"github.com/lm379/hmx/pkg/search"
	"github.com/meilisearch/meilisearch-go"
)

// SearchService 搜索服务
type SearchService struct {
	operaRepo  *repository.OperaRepo
	artistRepo *repository.ArtistRepo
}

// NewSearchService 创建搜索服务实例
func NewSearchService() *SearchService {
	return &SearchService{
		operaRepo:  repository.NewOperaRepo(),
		artistRepo: repository.NewArtistRepo(),
	}
}

// OperaSearchDocument Opera 搜索文档结构
type OperaSearchDocument struct {
	OperaID       uint     `json:"opera_id"`
	OperaTitle    string   `json:"opera_title"`
	OperaTitlePY  string   `json:"opera_title_py"` // 拼音字段
	ArtistNames   []string `json:"artist_names"`
	ArtistNamesPY []string `json:"artist_names_py"` // 艺术家名字拼音
	ArtistIDs     []uint   `json:"artist_ids"`
}

// ArtistSearchDocument Artist 搜索文档结构
type ArtistSearchDocument struct {
	ArtistID uint   `json:"artist_id"`
	Name     string `json:"name"`
	NamePY   string `json:"name_py"` // 拼音字段
}

// SearchOperasRequest 搜索 Opera 请求
type SearchOperasRequest struct {
	Query         string `json:"query" form:"query"`                   // 搜索关键词
	Page          int    `json:"page" form:"page"`                     // 页码
	PageSize      int    `json:"page_size" form:"page_size"`           // 每页数量
	IncludeHidden bool   `json:"include_hidden" form:"include_hidden"` // 是否包含隐藏内容
	SortBy        string `json:"sort_by" form:"sort_by"`               // 排序字段
}

// SearchOperasResponse 搜索 Opera 响应
type SearchOperasResponse struct {
	Results        []OperaSearchDocument `json:"results"`
	Total          int64                 `json:"total"`
	Page           int                   `json:"page"`
	PageSize       int                   `json:"page_size"`
	TotalPages     int                   `json:"total_pages"`
	ProcessingTime int64                 `json:"processing_time_ms"`
}

// SearchArtistsRequest 搜索 Artist 请求
type SearchArtistsRequest struct {
	Query    string `json:"query" form:"query"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	SortBy   string `json:"sort_by" form:"sort_by"`
}

// SearchArtistsResponse 搜索 Artist 响应
type SearchArtistsResponse struct {
	Results        []ArtistSearchDocument `json:"results"`
	Total          int64                  `json:"total"`
	Page           int                    `json:"page"`
	PageSize       int                    `json:"page_size"`
	TotalPages     int                    `json:"total_pages"`
	ProcessingTime int64                  `json:"processing_time_ms"`
}

// SearchOperas 搜索 Opera
func (s *SearchService) SearchOperas(req SearchOperasRequest) (*SearchOperasResponse, error) {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get search client: %w", err)
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	index := client.Index("operas")

	// 构建搜索请求
	searchReq := &meilisearch.SearchRequest{
		Limit:  int64(req.PageSize),
		Offset: int64((req.Page - 1) * req.PageSize),
	}

	// 添加排序
	if req.SortBy != "" {
		searchReq.Sort = []string{req.SortBy}
	}

	// 执行搜索
	searchResp, err := index.Search(req.Query, searchReq)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 解析结果
	results := make([]OperaSearchDocument, 0, len(searchResp.Hits))
	for _, hit := range searchResp.Hits {
		doc := OperaSearchDocument{}
		if err := hit.Decode(&doc); err != nil {
			log.Printf("Failed to decode search hit: %v", err)
			continue
		}
		results = append(results, doc)
	}

	totalPages := int(searchResp.EstimatedTotalHits) / req.PageSize
	if int(searchResp.EstimatedTotalHits)%req.PageSize != 0 {
		totalPages++
	}

	return &SearchOperasResponse{
		Results:        results,
		Total:          searchResp.EstimatedTotalHits,
		Page:           req.Page,
		PageSize:       req.PageSize,
		TotalPages:     totalPages,
		ProcessingTime: searchResp.ProcessingTimeMs,
	}, nil
}

// SearchArtists 搜索 Artist
func (s *SearchService) SearchArtists(req SearchArtistsRequest) (*SearchArtistsResponse, error) {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get search client: %w", err)
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	index := client.Index("artists")

	// 构建搜索请求
	searchReq := &meilisearch.SearchRequest{
		Limit:  int64(req.PageSize),
		Offset: int64((req.Page - 1) * req.PageSize),
	}

	// 添加排序
	if req.SortBy != "" {
		searchReq.Sort = []string{req.SortBy}
	}

	// 执行搜索
	searchResp, err := index.Search(req.Query, searchReq)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 解析结果
	results := make([]ArtistSearchDocument, 0, len(searchResp.Hits))
	for _, hit := range searchResp.Hits {
		doc := ArtistSearchDocument{}
		if err := hit.Decode(&doc); err != nil {
			log.Printf("Failed to decode search hit: %v", err)
			continue
		}
		results = append(results, doc)
	}

	totalPages := int(searchResp.EstimatedTotalHits) / req.PageSize
	if int(searchResp.EstimatedTotalHits)%req.PageSize != 0 {
		totalPages++
	}

	return &SearchArtistsResponse{
		Results:        results,
		Total:          searchResp.EstimatedTotalHits,
		Page:           req.Page,
		PageSize:       req.PageSize,
		TotalPages:     totalPages,
		ProcessingTime: searchResp.ProcessingTimeMs,
	}, nil
}

// IndexOpera 索引单个 Opera
func (s *SearchService) IndexOpera(opera *models.Opera) error {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("operas")

	// 构建搜索文档
	doc := s.operaToSearchDocument(opera)

	// 索引文档
	_, err = index.AddDocuments([]OperaSearchDocument{doc}, nil)
	if err != nil {
		return fmt.Errorf("failed to index opera: %w", err)
	}

	log.Printf("Indexed opera: %d - %s", opera.OperaID, opera.OperaTitle)
	return nil
}

// UpdateOperaIndex 更新 Opera 索引
func (s *SearchService) UpdateOperaIndex(opera *models.Opera) error {
	return s.IndexOpera(opera)
}

// DeleteOperaIndex 删除 Opera 索引
func (s *SearchService) DeleteOperaIndex(operaID uint) error {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("operas")
	_, err = index.DeleteDocument(strconv.FormatUint(uint64(operaID), 10), nil)
	if err != nil {
		return fmt.Errorf("failed to delete opera from index: %w", err)
	}

	log.Printf("Deleted opera from index: %d", operaID)
	return nil
}

// IndexArtist 索引单个 Artist
func (s *SearchService) IndexArtist(artist *models.Artist) error {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("artists")

	// 构建搜索文档
	doc := s.artistToSearchDocument(artist)

	// 索引文档
	_, err = index.AddDocuments([]ArtistSearchDocument{doc}, nil)
	if err != nil {
		return fmt.Errorf("failed to index artist: %w", err)
	}

	log.Printf("Indexed artist: %d - %s", artist.ArtistID, artist.Name)
	return nil
}

// UpdateArtistIndex 更新 Artist 索引
func (s *SearchService) UpdateArtistIndex(artist *models.Artist) error {
	return s.IndexArtist(artist)
}

// DeleteArtistIndex 删除 Artist 索引
func (s *SearchService) DeleteArtistIndex(artistID uint) error {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("artists")
	_, err = index.DeleteDocument(strconv.FormatUint(uint64(artistID), 10), nil)
	if err != nil {
		return fmt.Errorf("failed to delete artist from index: %w", err)
	}

	log.Printf("Deleted artist from index: %d", artistID)
	return nil
}

// ReindexAllOperas 重新索引所有 Opera
func (s *SearchService) ReindexAllOperas() error {
	// 获取所有 Opera（包括隐藏的）
	operas, err := s.operaRepo.GetAllForIndex()
	if err != nil {
		return fmt.Errorf("failed to get operas: %w", err)
	}

	if len(operas) == 0 {
		log.Println("No operas to index")
		return nil
	}

	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("operas")

	// 分批处理，每批50条
	batchSize := 50
	totalBatches := (len(operas) + batchSize - 1) / batchSize

	for i := 0; i < len(operas); i += batchSize {
		end := i + batchSize
		if end > len(operas) {
			end = len(operas)
		}

		batch := operas[i:end]
		docs := make([]OperaSearchDocument, len(batch))
		for j, opera := range batch {
			docs[j] = s.operaToSearchDocument(&opera)
		}

		_, err = index.AddDocuments(docs, nil)
		if err != nil {
			return fmt.Errorf("failed to batch index operas (batch %d/%d, range %d-%d): %w", (i/batchSize)+1, totalBatches, i, end-1, err)
		}

		log.Printf("Successfully indexed batch %d/%d (%d operas)", (i/batchSize)+1, totalBatches, len(batch))
	}

	log.Printf("Successfully reindexed %d operas in %d batches", len(operas), totalBatches)
	return nil
}

// ReindexAllArtists 重新索引所有 Artist
func (s *SearchService) ReindexAllArtists() error {
	// 获取所有 Artist
	artists, err := s.artistRepo.GetAllForIndex()
	if err != nil {
		return fmt.Errorf("failed to get artists: %w", err)
	}

	if len(artists) == 0 {
		log.Println("No artists to index")
		return nil
	}

	client, err := search.GetMeilisearchClient()
	if err != nil {
		return fmt.Errorf("failed to get search client: %w", err)
	}

	index := client.Index("artists")

	// 分批处理，每批50条
	batchSize := 50
	totalBatches := (len(artists) + batchSize - 1) / batchSize

	for i := 0; i < len(artists); i += batchSize {
		end := i + batchSize
		if end > len(artists) {
			end = len(artists)
		}

		batch := artists[i:end]
		docs := make([]ArtistSearchDocument, len(batch))
		for j, artist := range batch {
			docs[j] = s.artistToSearchDocument(&artist)
		}

		_, err = index.AddDocuments(docs, nil)
		if err != nil {
			return fmt.Errorf("failed to batch index artists (batch %d/%d, range %d-%d): %w", (i/batchSize)+1, totalBatches, i, end-1, err)
		}

		log.Printf("Successfully indexed batch %d/%d (%d artists)", (i/batchSize)+1, totalBatches, len(batch))
	}

	log.Printf("Successfully reindexed %d artists in %d batches", len(artists), totalBatches)
	return nil
}

// operaToSearchDocument 将 Opera 模型转换为搜索文档
func (s *SearchService) operaToSearchDocument(opera *models.Opera) OperaSearchDocument {
	doc := OperaSearchDocument{
		OperaID:      opera.OperaID,
		OperaTitle:   opera.OperaTitle,
		OperaTitlePY: pinyin.ToPinyin(opera.OperaTitle), // 添加拼音
	}

	// 提取艺术家信息
	if len(opera.Artists) > 0 {
		doc.ArtistNames = make([]string, len(opera.Artists))
		doc.ArtistNamesPY = make([]string, len(opera.Artists))
		doc.ArtistIDs = make([]uint, len(opera.Artists))
		for i, artist := range opera.Artists {
			doc.ArtistNames[i] = artist.Name
			doc.ArtistNamesPY[i] = pinyin.ToPinyin(artist.Name) // 添加拼音
			doc.ArtistIDs[i] = artist.ArtistID
		}
	}

	return doc
}

// artistToSearchDocument 将 Artist 模型转换为搜索文档
func (s *SearchService) artistToSearchDocument(artist *models.Artist) ArtistSearchDocument {
	return ArtistSearchDocument{
		ArtistID: artist.ArtistID,
		Name:     artist.Name,
		NamePY:   pinyin.ToPinyin(artist.Name), // 添加拼音
	}
}

// GetSearchStats 获取搜索统计
func (s *SearchService) GetSearchStats() (map[string]interface{}, error) {
	client, err := search.GetMeilisearchClient()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})

	// 获取 operas 索引统计
	operaIndex := client.Index("operas")
	operaStats, err := operaIndex.GetStats()
	if err == nil {
		stats["operas"] = operaStats
	}

	// 获取 artists 索引统计
	artistIndex := client.Index("artists")
	artistStats, err := artistIndex.GetStats()
	if err == nil {
		stats["artists"] = artistStats
	}

	return stats, nil
}

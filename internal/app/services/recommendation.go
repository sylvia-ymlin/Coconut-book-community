package services

import (
	"github.com/yourusername/bookcommunity/internal/app/models"
	// "bytes"
	// "context"
	// "encoding/json"
	// "fmt"
	// "net/http"
	// "time"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	// 预留：未来可配置Python推荐API地址
	// pythonAPIUrl string
	// httpClient   *http.Client
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{
		// pythonAPIUrl: "http://localhost:6006",
		// httpClient: &http.Client{
		// 	Timeout: 3 * time.Second,
		// },
	}
}

// GetPersonalizedRecommendations 获取个性化推荐
// TODO: 未来可对接真实推荐系统
// 当前返回mock数据
func (s *RecommendationService) GetPersonalizedRecommendations(userID uint, topK int) ([]*models.Book, error) {
	// 检查配置：是否启用真实推荐
	// if config.GetRecommendConfig().Enabled {
	// 	return s.getRemoteRecommendations(userID, topK)
	// }

	// 当前返回mock数据
	return s.getMockRecommendations(userID, topK), nil
}

// SemanticSearch 语义搜索
// TODO: 未来可对接RAG检索系统
// 当前返回mock数据
func (s *RecommendationService) SemanticSearch(query string, topK int) ([]*models.Book, error) {
	// if config.GetRecommendConfig().Enabled {
	// 	return s.getRemoteSearch(query, topK)
	// }

	// 当前返回mock搜索结果
	return s.getMockSearchResults(query, topK), nil
}

// getMockRecommendations 生成mock推荐数据
func (s *RecommendationService) getMockRecommendations(userID uint, topK int) []*models.Book {
	mockBooks := []*models.Book{
		{
			ISBN:      "9787111544937",
			Title:     "深入理解计算机系统（原书第3版）",
			Author:    "Randal E. Bryant / David R. O'Hallaron",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s29195878.jpg",
			Rating:    9.7,
			Reason:    "基于你的阅读历史推荐",
			Publisher: "机械工业出版社",
			PubDate:   "2016-11",
			Summary:   "从程序员的视角，看计算机系统！本书适用于那些想要写出更快、更可靠程序的程序员。",
		},
		{
			ISBN:      "9787115428028",
			Title:     "Go语言圣经",
			Author:    "Alan A. A. Donovan / Brian W. Kernighan",
			CoverURL:  "https://img9.doubanio.com/view/subject/l/public/s28699046.jpg",
			Rating:    9.5,
			Reason:    "编程类畅销书",
			Publisher: "人民邮电出版社",
			PubDate:   "2017-1",
			Summary:   "Go语言圣经是Go语言领域最经典、最权威的书籍之一。",
		},
		{
			ISBN:      "9787111421900",
			Title:     "编码：隐匿在计算机软硬件背后的语言",
			Author:    "Charles Petzold",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s26490404.jpg",
			Rating:    9.3,
			Reason:    "计算机科学经典",
			Publisher: "机械工业出版社",
			PubDate:   "2012-10",
			Summary:   "一本讲述计算机工作原理的书。从电灯开关的原理开始，到计算机的组成，层层递进。",
		},
		{
			ISBN:      "9787111213826",
			Title:     "代码大全（第2版）",
			Author:    "Steve McConnell",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s1495029.jpg",
			Rating:    9.3,
			Reason:    "软件工程必读",
			Publisher: "电子工业出版社",
			PubDate:   "2006-3",
			Summary:   "软件开发领域的百科全书，无论你是初学者还是资深开发者，都能从中受益。",
		},
		{
			ISBN:      "9787115385130",
			Title:     "算法（第4版）",
			Author:    "Robert Sedgewick / Kevin Wayne",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s28322244.jpg",
			Rating:    9.4,
			Reason:    "算法经典教材",
			Publisher: "人民邮电出版社",
			PubDate:   "2012-10",
			Summary:   "全面介绍了算法和数据结构的必备知识，包含排序、搜索、图处理和字符串处理等。",
		},
		{
			ISBN:      "9787115291028",
			Title:     "计算机程序的构造和解释（原书第2版）",
			Author:    "Harold Abelson / Gerald Jay Sussman",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s6183487.jpg",
			Rating:    9.5,
			Reason:    "编程思想经典",
			Publisher: "机械工业出版社",
			PubDate:   "2004-2",
			Summary:   "SICP是MIT的经典教材，教授编程的本质，而不仅仅是语法。",
		},
		{
			ISBN:      "9787115275790",
			Title:     "设计模式：可复用面向对象软件的基础",
			Author:    "Erich Gamma / Richard Helm / Ralph Johnson / John Vlissides",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s1074361.jpg",
			Rating:    9.1,
			Reason:    "设计模式权威指南",
			Publisher: "机械工业出版社",
			PubDate:   "2000-9",
			Summary:   "四人帮的经典著作，介绍了23种设计模式。",
		},
		{
			ISBN:      "9787115385390",
			Title:     "数据结构与算法分析：C语言描述（原书第2版）",
			Author:    "Mark Allen Weiss",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s28299194.jpg",
			Rating:    9.0,
			Reason:    "数据结构经典",
			Publisher: "机械工业出版社",
			PubDate:   "2004-1",
			Summary:   "系统介绍数据结构和算法分析的经典教材。",
		},
		{
			ISBN:      "9787115449689",
			Title:     "Python编程：从入门到实践",
			Author:    "Eric Matthes",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s28891804.jpg",
			Rating:    9.1,
			Reason:    "Python入门首选",
			Publisher: "人民邮电出版社",
			PubDate:   "2016-7",
			Summary:   "一本全面的Python编程从入门到实践教程。",
		},
		{
			ISBN:      "9787115373991",
			Title:     "Effective Java中文版（第2版）",
			Author:    "Joshua Bloch",
			CoverURL:  "https://img3.doubanio.com/view/subject/l/public/s27243455.jpg",
			Rating:    9.1,
			Reason:    "Java最佳实践",
			Publisher: "机械工业出版社",
			PubDate:   "2009-1",
			Summary:   "Java编程的最佳实践指南，包含78条编程准则。",
		},
	}

	// 限制返回数量
	if topK > 0 && topK < len(mockBooks) {
		return mockBooks[:topK]
	}
	return mockBooks
}

// getMockSearchResults 生成mock搜索结果
func (s *RecommendationService) getMockSearchResults(query string, topK int) []*models.Book {
	// 简单的关键词匹配逻辑
	// 当前直接返回所有mock数据的前N本
	allBooks := s.getMockRecommendations(0, 20)

	// TODO: 实现简单的关键词过滤
	// 例如：如果query包含"Go"，优先返回Go相关书籍
	// 当前简化处理：直接返回前topK本

	if topK > 0 && topK < len(allBooks) {
		return allBooks[:topK]
	}
	return allBooks
}

// getRemoteRecommendations 调用真实推荐API（预留）
// 未来可以这样实现：
//
// func (s *RecommendationService) getRemoteRecommendations(userID uint, topK int) ([]*models.Book, error) {
// 	reqBody := map[string]interface{}{
// 		"user_id": userID,
// 		"top_k":   topK,
// 	}
//
// 	body, err := json.Marshal(reqBody)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
//
// 	req, err := http.NewRequestWithContext(ctx, "POST",
// 		s.pythonAPIUrl+"/api/v1/recommend/personalized",
// 		bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
//
// 	resp, err := s.httpClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("request failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	var result struct {
// 		StatusCode int            `json:"status_code"`
// 		Books      []*models.Book `json:"books"`
// 		Message    string         `json:"message"`
// 	}
//
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("decode failed: %w", err)
// 	}
//
// 	if result.StatusCode != 0 {
// 		return nil, fmt.Errorf("recommendation failed: %s", result.Message)
// 	}
//
// 	return result.Books, nil
// }

// getRemoteSearch 调用真实搜索API（预留）
//
// func (s *RecommendationService) getRemoteSearch(query string, topK int) ([]*models.Book, error) {
// 	reqBody := map[string]interface{}{
// 		"query": query,
// 		"top_k": topK,
// 	}
//
// 	body, _ := json.Marshal(reqBody)
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
//
// 	req, _ := http.NewRequestWithContext(ctx, "POST",
// 		s.pythonAPIUrl+"/api/v1/search/semantic",
// 		bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
//
// 	resp, err := s.httpClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("request failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	var result struct {
// 		StatusCode int            `json:"status_code"`
// 		Books      []*models.Book `json:"books"`
// 	}
//
// 	json.NewDecoder(resp.Body).Decode(&result)
// 	return result.Books, nil
// }

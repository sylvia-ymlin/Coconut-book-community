package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/sylvia-ymlin/Coconut-book-community/internal/app/models"
	"github.com/sylvia-ymlin/Coconut-book-community/pkg/utils"
)

var logger = utils.NewLogger("recommendation_client")

// RecommendationClient HTTP 客户端，用于调用 Python 推荐服务
type RecommendationClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewRecommendationClient 创建推荐服务客户端
func NewRecommendationClient(baseURL string) *RecommendationClient {
	return &RecommendationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ==================== 响应结构体 ====================

// SearchResponse Python 搜索服务响应
type SearchResponse struct {
	Results []models.Book `json:"results"`
	Total   int           `json:"total"`
}

// RecommendResponse Python 推荐服务响应
type RecommendResponse struct {
	Books []models.Book `json:"books"`
	Total int           `json:"total"`
}

// BookDetailResponse Python 图书详情响应
type BookDetailResponse struct {
	Book models.Book `json:"book"`
}

// ChatRequest Chat with Book 请求
type ChatRequest struct {
	ISBN    string `json:"isbn"`
	Message string `json:"message"`
	UserID  uint   `json:"user_id,omitempty"`
}

// ChatResponse Chat with Book 响应
type ChatResponse struct {
	Response string `json:"response"`
	ISBN     string `json:"isbn"`
}

// ==================== API 方法 ====================

// SearchBooks 调用 Python RAG 搜索服务
// 参数:
//   - query: 搜索关键词（支持语义搜索）
//   - topK: 返回结果数量
// 返回: 图书列表和错误信息
func (r *RecommendationClient) SearchBooks(query string, topK int) ([]models.Book, error) {
	// 构建 URL
	apiURL := fmt.Sprintf("%s/search", r.baseURL)
	params := url.Values{}
	params.Add("q", query)
	params.Add("top_k", fmt.Sprintf("%d", topK))

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	logger.Printf("Calling Python search API: %s", fullURL)

	// 发送 HTTP GET 请求
	resp, err := r.httpClient.Get(fullURL)
	if err != nil {
		logger.Printf("Failed to call search API: %v", err)
		return nil, fmt.Errorf("failed to call search API: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Printf("Search API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("search API returned status %d", resp.StatusCode)
	}

	// 解析响应
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		logger.Printf("Failed to decode search response: %v", err)
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	logger.Printf("Search API returned %d books", len(searchResp.Results))
	return searchResp.Results, nil
}

// GetPersonalRecommendations 调用 Python 个性化推荐服务
// 参数:
//   - userID: 用户ID
//   - topK: 返回结果数量
// 返回: 推荐图书列表和错误信息
func (r *RecommendationClient) GetPersonalRecommendations(userID uint, topK int) ([]models.Book, error) {
	// 构建 URL
	apiURL := fmt.Sprintf("%s/recommend/personal", r.baseURL)
	params := url.Values{}
	params.Add("user_id", fmt.Sprintf("%d", userID))
	params.Add("top_k", fmt.Sprintf("%d", topK))

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	logger.Printf("Calling Python recommendation API: %s", fullURL)

	// 发送 HTTP GET 请求
	resp, err := r.httpClient.Get(fullURL)
	if err != nil {
		logger.Printf("Failed to call recommendation API: %v", err)
		return nil, fmt.Errorf("failed to call recommendation API: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Printf("Recommendation API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("recommendation API returned status %d", resp.StatusCode)
	}

	// 解析响应
	var recommendResp RecommendResponse
	if err := json.NewDecoder(resp.Body).Decode(&recommendResp); err != nil {
		logger.Printf("Failed to decode recommendation response: %v", err)
		return nil, fmt.Errorf("failed to decode recommendation response: %w", err)
	}

	logger.Printf("Recommendation API returned %d books", len(recommendResp.Books))
	return recommendResp.Books, nil
}

// GetBookDetail 获取图书详情
// 参数:
//   - isbn: 图书 ISBN
// 返回: 图书详情和错误信息
func (r *RecommendationClient) GetBookDetail(isbn string) (*models.Book, error) {
	// 构建 URL
	apiURL := fmt.Sprintf("%s/books/%s", r.baseURL, isbn)
	logger.Printf("Calling Python book detail API: %s", apiURL)

	// 发送 HTTP GET 请求
	resp, err := r.httpClient.Get(apiURL)
	if err != nil {
		logger.Printf("Failed to call book detail API: %v", err)
		return nil, fmt.Errorf("failed to call book detail API: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Printf("Book detail API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("book detail API returned status %d", resp.StatusCode)
	}

	// 解析响应
	var bookResp BookDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookResp); err != nil {
		logger.Printf("Failed to decode book detail response: %v", err)
		return nil, fmt.Errorf("failed to decode book detail response: %w", err)
	}

	logger.Printf("Book detail API returned book: %s", bookResp.Book.Title)
	return &bookResp.Book, nil
}

// ChatWithBook 与图书进行对话（调用 Python LLM 服务）
// 参数:
//   - isbn: 图书 ISBN
//   - message: 用户消息
//   - userID: 用户ID（可选）
// 返回: LLM 响应和错误信息
func (r *RecommendationClient) ChatWithBook(isbn, message string, userID uint) (string, error) {
	// 构建 URL
	apiURL := fmt.Sprintf("%s/chat", r.baseURL)
	logger.Printf("Calling Python chat API for ISBN: %s", isbn)

	// 构建请求体
	chatReq := ChatRequest{
		ISBN:    isbn,
		Message: message,
		UserID:  userID,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		logger.Printf("Failed to marshal chat request: %v", err)
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	// 发送 HTTP POST 请求
	resp, err := r.httpClient.Post(apiURL, "application/json", nil)
	if err != nil {
		logger.Printf("Failed to call chat API: %v", err)
		return "", fmt.Errorf("failed to call chat API: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Printf("Chat API returned status %d: %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("chat API returned status %d", resp.StatusCode)
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		logger.Printf("Failed to decode chat response: %v", err)
		return "", fmt.Errorf("failed to decode chat response: %w", err)
	}

	logger.Printf("Chat API returned response for ISBN: %s", isbn)
	return chatResp.Response, nil
}

// HealthCheck 检查 Python 服务健康状态
func (r *RecommendationClient) HealthCheck() error {
	apiURL := fmt.Sprintf("%s/health", r.baseURL)
	logger.Printf("Checking Python service health: %s", apiURL)

	resp, err := r.httpClient.Get(apiURL)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	logger.Printf("Python service is healthy")
	return nil
}

// ==================== 全局客户端实例 ====================

var (
	// GlobalRecommendationClient 全局推荐服务客户端
	// 在 main.go 中初始化
	GlobalRecommendationClient *RecommendationClient
)

// InitRecommendationClient 初始化全局推荐服务客户端
func InitRecommendationClient(baseURL string) {
	GlobalRecommendationClient = NewRecommendationClient(baseURL)
	logger.Printf("Initialized recommendation client with base URL: %s", baseURL)

	// 健康检查
	if err := GlobalRecommendationClient.HealthCheck(); err != nil {
		logger.Printf("Warning: Python recommendation service health check failed: %v", err)
	}
}

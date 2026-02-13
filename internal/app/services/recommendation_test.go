package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecommendationService(t *testing.T) {
	service := NewRecommendationService()
	assert.NotNil(t, service)
}

func TestGetPersonalizedRecommendations(t *testing.T) {
	service := NewRecommendationService()

	t.Run("get recommendations with default top_k", func(t *testing.T) {
		userID := uint(1)
		topK := 10

		books, err := service.GetPersonalizedRecommendations(userID, topK)

		assert.NoError(t, err)
		assert.NotNil(t, books)
		assert.LessOrEqual(t, len(books), topK)
		assert.Greater(t, len(books), 0)

		// Verify first book structure
		if len(books) > 0 {
			book := books[0]
			assert.NotEmpty(t, book.ISBN)
			assert.NotEmpty(t, book.Title)
			assert.NotEmpty(t, book.Author)
			assert.Greater(t, book.Rating, 0.0)
			assert.NotEmpty(t, book.Reason)
		}
	})

	t.Run("get recommendations with small top_k", func(t *testing.T) {
		userID := uint(2)
		topK := 3

		books, err := service.GetPersonalizedRecommendations(userID, topK)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(books), topK)
		assert.Greater(t, len(books), 0)
	})

	t.Run("get recommendations with large top_k", func(t *testing.T) {
		userID := uint(3)
		topK := 50

		books, err := service.GetPersonalizedRecommendations(userID, topK)

		assert.NoError(t, err)
		assert.NotNil(t, books)
		// Mock data might return less than requested
		assert.Greater(t, len(books), 0)
	})

	t.Run("different users get same mock data", func(t *testing.T) {
		// In mock mode, different users get similar recommendations
		books1, _ := service.GetPersonalizedRecommendations(1, 5)
		books2, _ := service.GetPersonalizedRecommendations(2, 5)

		assert.Equal(t, len(books1), len(books2))
		// Mock data should have consistent structure
		if len(books1) > 0 && len(books2) > 0 {
			assert.Equal(t, books1[0].ISBN, books2[0].ISBN)
		}
	})
}

func TestSemanticSearch(t *testing.T) {
	service := NewRecommendationService()

	t.Run("search with valid query", func(t *testing.T) {
		query := "golang"
		topK := 10

		books, err := service.SemanticSearch(query, topK)

		assert.NoError(t, err)
		assert.NotNil(t, books)
		assert.LessOrEqual(t, len(books), topK)
		assert.Greater(t, len(books), 0)

		// Verify book structure
		if len(books) > 0 {
			book := books[0]
			assert.NotEmpty(t, book.ISBN)
			assert.NotEmpty(t, book.Title)
			assert.NotEmpty(t, book.Author)
			assert.Greater(t, book.Rating, 0.0)
		}
	})

	t.Run("search with chinese query", func(t *testing.T) {
		query := "计算机"
		topK := 5

		books, err := service.SemanticSearch(query, topK)

		assert.NoError(t, err)
		assert.Greater(t, len(books), 0)
	})

	t.Run("search with empty query", func(t *testing.T) {
		query := ""
		topK := 10

		books, err := service.SemanticSearch(query, topK)

		// Mock implementation should still return results
		assert.NoError(t, err)
		assert.NotNil(t, books)
	})

	t.Run("search with small top_k", func(t *testing.T) {
		query := "programming"
		topK := 2

		books, err := service.SemanticSearch(query, topK)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(books), topK)
	})
}

func TestMockRecommendations(t *testing.T) {
	service := NewRecommendationService()

	t.Run("mock data has valid structure", func(t *testing.T) {
		books := service.getMockRecommendations(1, 5)

		assert.Greater(t, len(books), 0)

		for _, book := range books {
			// ISBN should be valid (13 digits)
			assert.Len(t, book.ISBN, 13)
			assert.NotEmpty(t, book.Title)
			assert.NotEmpty(t, book.Author)
			assert.Greater(t, book.Rating, 0.0)
			assert.LessOrEqual(t, book.Rating, 10.0)
			assert.NotEmpty(t, book.Reason)
			assert.NotEmpty(t, book.Publisher)
		}
	})

	t.Run("respects top_k limit", func(t *testing.T) {
		topK := 3
		books := service.getMockRecommendations(1, topK)

		assert.LessOrEqual(t, len(books), topK)
	})
}

func TestMockSearchResults(t *testing.T) {
	service := NewRecommendationService()

	t.Run("mock search returns results", func(t *testing.T) {
		books := service.getMockSearchResults("test", 5)

		assert.Greater(t, len(books), 0)

		for _, book := range books {
			assert.NotEmpty(t, book.ISBN)
			assert.NotEmpty(t, book.Title)
			assert.Greater(t, book.Rating, 0.0)
		}
	})

	t.Run("search respects top_k", func(t *testing.T) {
		topK := 2
		books := service.getMockSearchResults("golang", topK)

		assert.LessOrEqual(t, len(books), topK)
	})
}

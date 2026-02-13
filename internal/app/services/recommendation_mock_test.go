package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the mock data generation functions independently
func TestGetMockRecommendations(t *testing.T) {
	service := NewRecommendationService()

	t.Run("returns valid mock data", func(t *testing.T) {
		books := service.getMockRecommendations(1, 5)
		assert.NotNil(t, books)
		assert.Greater(t, len(books), 0)
	})

	t.Run("respects topK parameter", func(t *testing.T) {
		topK := 3
		books := service.getMockRecommendations(1, topK)
		assert.LessOrEqual(t, len(books), topK)
	})
}

func TestGetMockSearchResults(t *testing.T) {
	service := NewRecommendationService()

	t.Run("returns search results", func(t *testing.T) {
		books := service.getMockSearchResults("test", 5)
		assert.NotNil(t, books)
		assert.Greater(t, len(books), 0)
	})

	t.Run("respects topK", func(t *testing.T) {
		topK := 2
		books := service.getMockSearchResults("golang", topK)
		assert.LessOrEqual(t, len(books), topK)
	})
}

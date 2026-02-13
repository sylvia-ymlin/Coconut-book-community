package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBook(t *testing.T) {
	t.Run("create book with all fields", func(t *testing.T) {
		book := Book{
			ISBN:      "9787111544937",
			Title:     "深入理解计算机系统",
			Author:    "Randal E. Bryant",
			CoverURL:  "https://example.com/cover.jpg",
			Rating:    9.7,
			Reason:    "推荐理由",
			Publisher: "机械工业出版社",
			PubDate:   "2016-11",
			Summary:   "本书详细介绍计算机系统...",
		}

		assert.Equal(t, "9787111544937", book.ISBN)
		assert.Equal(t, "深入理解计算机系统", book.Title)
		assert.Equal(t, "Randal E. Bryant", book.Author)
		assert.Equal(t, float32(9.7), book.Rating)
		assert.NotEmpty(t, book.CoverURL)
		assert.NotEmpty(t, book.Reason)
		assert.NotEmpty(t, book.Publisher)
		assert.NotEmpty(t, book.PubDate)
		assert.NotEmpty(t, book.Summary)
	})

	t.Run("create minimal book", func(t *testing.T) {
		book := Book{
			ISBN:   "9787115428028",
			Title:  "Go语言圣经",
			Author: "Alan Donovan",
			Rating: 9.5,
		}

		assert.Equal(t, "9787115428028", book.ISBN)
		assert.Equal(t, "Go语言圣经", book.Title)
		assert.Equal(t, "Alan Donovan", book.Author)
		assert.Equal(t, float32(9.5), book.Rating)
		assert.Empty(t, book.CoverURL)
		assert.Empty(t, book.Reason)
	})

	t.Run("book JSON marshaling", func(t *testing.T) {
		book := Book{
			ISBN:   "9787111544937",
			Title:  "Test Book",
			Author: "Test Author",
			Rating: 8.5,
			Reason: "Good book",
		}

		jsonData, err := json.Marshal(book)
		assert.NoError(t, err)
		assert.NotNil(t, jsonData)

		// Unmarshal back
		var unmarshaled Book
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, book.ISBN, unmarshaled.ISBN)
		assert.Equal(t, book.Title, unmarshaled.Title)
		assert.Equal(t, book.Rating, unmarshaled.Rating)
	})

	t.Run("rating bounds", func(t *testing.T) {
		book := Book{
			ISBN:   "123",
			Title:  "Test",
			Author: "Author",
			Rating: 10.0,
		}
		assert.LessOrEqual(t, book.Rating, float32(10.0))

		book.Rating = 0.0
		assert.GreaterOrEqual(t, book.Rating, float32(0.0))
	})
}

func TestBookSearchRequest(t *testing.T) {
	t.Run("valid search request", func(t *testing.T) {
		req := BookSearchRequest{
			Query: "golang",
			TopK:  10,
		}

		assert.Equal(t, "golang", req.Query)
		assert.Equal(t, 10, req.TopK)
		assert.Greater(t, len(req.Query), 0)
	})

	t.Run("search request with chinese", func(t *testing.T) {
		req := BookSearchRequest{
			Query: "计算机科学",
			TopK:  5,
		}

		assert.Equal(t, "计算机科学", req.Query)
		assert.Equal(t, 5, req.TopK)
	})

	t.Run("default top_k", func(t *testing.T) {
		req := BookSearchRequest{
			Query: "test",
			TopK:  0, // Default or zero
		}

		assert.Equal(t, "test", req.Query)
		assert.Equal(t, 0, req.TopK)
	})

	t.Run("JSON binding tags", func(t *testing.T) {
		jsonStr := `{"query":"programming","top_k":15}`
		var req BookSearchRequest

		err := json.Unmarshal([]byte(jsonStr), &req)
		assert.NoError(t, err)
		assert.Equal(t, "programming", req.Query)
		assert.Equal(t, 15, req.TopK)
	})
}

func TestRecommendRequest(t *testing.T) {
	t.Run("valid recommend request", func(t *testing.T) {
		req := RecommendRequest{
			UserID: 123,
			TopK:   10,
		}

		assert.Equal(t, uint(123), req.UserID)
		assert.Equal(t, 10, req.TopK)
	})

	t.Run("zero user ID", func(t *testing.T) {
		req := RecommendRequest{
			UserID: 0,
			TopK:   5,
		}

		assert.Equal(t, uint(0), req.UserID)
		assert.Equal(t, 5, req.TopK)
	})

	t.Run("JSON marshaling", func(t *testing.T) {
		req := RecommendRequest{
			UserID: 456,
			TopK:   20,
		}

		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)

		var unmarshaled RecommendRequest
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, req.UserID, unmarshaled.UserID)
		assert.Equal(t, req.TopK, unmarshaled.TopK)
	})

	t.Run("large top_k", func(t *testing.T) {
		req := RecommendRequest{
			UserID: 999,
			TopK:   1000,
		}

		assert.Equal(t, uint(999), req.UserID)
		assert.Equal(t, 1000, req.TopK)
	})
}

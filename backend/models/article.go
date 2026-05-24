package models

import "time"

type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Content     string    `json:"content"`
	ScrapedAt   time.Time `json:"scrapedAt"`
	WordCount   int       `json:"wordCount"`
	Description string    `json:"description"`
}

type ArticleSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	WordCount   int    `json:"wordCount"`
	Description string `json:"description"`
}

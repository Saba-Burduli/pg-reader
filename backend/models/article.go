package models

import "time"

type Article struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	URL                 string    `json:"url"`
	Content             string    `json:"content"`
	ScrapedAt           time.Time `json:"scrapedAt"`
	PublishedAt         time.Time `json:"publishedAt"`
	PublishedDateSource string    `json:"publishedDateSource"`
	WordCount           int       `json:"wordCount"`
	Description         string    `json:"description"`
	IsRead              bool      `json:"isRead"`
}

type ArticleSummary struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	URL                 string    `json:"url"`
	PublishedAt         time.Time `json:"publishedAt"`
	PublishedDateSource string    `json:"publishedDateSource"`
	WordCount           int       `json:"wordCount"`
	Description         string    `json:"description"`
	IsRead              bool      `json:"isRead"`
}

package services

import (
	"context"
	"errors"

	"pgreader/models"
)

var ErrNotFound = errors.New("article not found")

type ArticleService struct {
	store   *Store
	scraper *Scraper
}

func NewArticleService(store *Store, scraper *Scraper) *ArticleService {
	return &ArticleService{store: store, scraper: scraper}
}

func (s *ArticleService) EnsureSynced(ctx context.Context) error {
	if !s.store.Empty() {
		return nil
	}
	return s.Sync(ctx)
}

func (s *ArticleService) Sync(ctx context.Context) error {
	articles, err := s.scraper.Sync(ctx)
	if err != nil {
		return err
	}
	return s.store.SaveAll(articles)
}

func (s *ArticleService) List() []models.ArticleSummary {
	return s.store.List()
}

func (s *ArticleService) Get(id string) (models.Article, error) {
	article, ok := s.store.Get(id)
	if !ok {
		return models.Article{}, ErrNotFound
	}
	return article, nil
}

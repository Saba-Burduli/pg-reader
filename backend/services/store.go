package services

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"pgreader/models"
)

type Store struct {
	mu       sync.RWMutex
	dataPath string
	articles map[string]models.Article
}

func NewStore(dataPath string) (*Store, error) {
	s := &Store{
		dataPath: dataPath,
		articles: make(map[string]models.Article),
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.dataPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	var articles []models.Article
	if err := json.NewDecoder(file).Decode(&articles); err != nil {
		return err
	}

	for _, a := range articles {
		s.articles[a.ID] = a
	}

	return nil
}

func (s *Store) SaveAll(articles []models.Article) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})
	if err := s.persistLocked(articles); err != nil {
		return err
	}

	s.articles = make(map[string]models.Article, len(articles))
	for _, a := range articles {
		s.articles[a.ID] = a
	}

	return nil
}

func (s *Store) List() []models.ArticleSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.ArticleSummary, 0, len(s.articles))
	for _, a := range s.articles {
		out = append(out, models.ArticleSummary{
			ID:                  a.ID,
			Title:               a.Title,
			URL:                 a.URL,
			PublishedAt:         a.PublishedAt,
			PublishedDateSource: a.PublishedDateSource,
			WordCount:           a.WordCount,
			Description:         a.Description,
			IsRead:              a.IsRead,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].PublishedAt.After(out[j].PublishedAt)
	})

	return out
}

func (s *Store) Get(id string) (models.Article, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.articles[id]
	return a, ok
}

func (s *Store) Empty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.articles) == 0
}

func (s *Store) NeedsMetadataRefresh() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.articles) == 0 {
		return true
	}
	for _, article := range s.articles {
		if article.PublishedAt.IsZero() {
			return true
		}
	}
	return false
}

func (s *Store) ReadStateMap() map[string]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]bool, len(s.articles))
	for id, article := range s.articles {
		out[id] = article.IsRead
	}
	return out
}

func (s *Store) SetRead(id string, isRead bool) (models.Article, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	article, ok := s.articles[id]
	if !ok {
		return models.Article{}, false, nil
	}
	article.IsRead = isRead
	s.articles[id] = article

	articles := make([]models.Article, 0, len(s.articles))
	for _, a := range s.articles {
		articles = append(articles, a)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt.After(articles[j].PublishedAt)
	})

	if err := s.persistLocked(articles); err != nil {
		return models.Article{}, false, err
	}

	return article, true, nil
}

func (s *Store) persistLocked(articles []models.Article) error {
	if err := os.MkdirAll(filepath.Dir(s.dataPath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.dataPath, data, 0o644); err != nil {
		return err
	}
	return nil
}

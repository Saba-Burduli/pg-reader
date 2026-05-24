package services

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"

	"pgreader/models"
)

const (
	indexURL = "https://www.paulgraham.com/articles.html"
)

type Scraper struct {
	client      *http.Client
	rateLimiter <-chan time.Time
	retries     int
}

var (
	monthDayYearPattern = regexp.MustCompile(`(?i)\b(january|february|march|april|may|june|july|august|september|october|november|december)\s+\d{1,2},\s+\d{4}\b`)
	monthYearPattern    = regexp.MustCompile(`(?i)\b(january|february|march|april|may|june|july|august|september|october|november|december)\s+\d{4}\b`)
	shortDatePattern    = regexp.MustCompile(`\b\d{1,2}/\d{1,2}/\d{2,4}\b`)
)

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		rateLimiter: time.Tick(400 * time.Millisecond),
		retries:     3,
	}
}

func (s *Scraper) Sync(ctx context.Context) ([]models.Article, error) {
	doc, err := s.fetchDoc(ctx, indexURL)
	if err != nil {
		return nil, fmt.Errorf("fetch index: %w", err)
	}

	type entry struct {
		title string
		url   string
	}
	entries := make([]entry, 0, 200)
	seen := map[string]struct{}{}

	doc.Find("a").Each(func(_ int, sel *goquery.Selection) {
		href, ok := sel.Attr("href")
		if !ok {
			return
		}
		title := strings.TrimSpace(sel.Text())
		if title == "" {
			return
		}
		if !strings.HasSuffix(strings.ToLower(href), ".html") {
			return
		}
		if strings.Contains(strings.ToLower(href), "index") || strings.Contains(strings.ToLower(href), "rss") {
			return
		}
		url := "https://www.paulgraham.com/" + strings.TrimPrefix(href, "/")
		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}
		entries = append(entries, entry{title: title, url: url})
	})

	articles := make([]models.Article, 0, len(entries))
	for _, e := range entries {
		content, publishedAt, source, err := s.fetchArticleContent(ctx, e.url)
		if err != nil {
			continue
		}
		content = normalizeText(content)
		if content == "" {
			continue
		}
		id := slugify(e.title)
		articles = append(articles, models.Article{
			ID:                  id,
			Title:               e.title,
			URL:                 e.url,
			Content:             content,
			ScrapedAt:           time.Now().UTC(),
			PublishedAt:         publishedAt,
			PublishedDateSource: source,
			WordCount:           countWords(content),
			Description:         summarize(content),
			IsRead:              false,
		})
	}

	return articles, nil
}

func (s *Scraper) fetchArticleContent(ctx context.Context, url string) (string, time.Time, string, error) {
	doc, err := s.fetchDoc(ctx, url)
	if err != nil {
		return "", time.Time{}, "", err
	}
	doc.Find("script, style, img, noscript").Remove()

	var lineBuilder strings.Builder
	doc.Find("body").Children().Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "" {
			return
		}
		lineBuilder.WriteString(text)
		lineBuilder.WriteString("\n")
	})
	lineText := strings.TrimSpace(lineBuilder.String())
	if lineText == "" {
		lineText = strings.TrimSpace(doc.Find("body").Text())
	}

	publishedAt, source := inferPublishedDate(doc, lineText)
	return lineText, publishedAt, source, nil
}

func (s *Scraper) fetchDoc(ctx context.Context, url string) (*goquery.Document, error) {
	var lastErr error
	for i := 0; i < s.retries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-s.rateLimiter:
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "pg-reader/1.0 (local study project)")

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
			continue
		}
		if resp.StatusCode >= 400 {
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
			time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
			continue
		}
		return doc, nil
	}
	return nil, lastErr
}

func normalizeText(in string) string {
	var b strings.Builder
	b.Grow(len(in))

	prevSpace := false
	for _, r := range in {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash {
			b.WriteRune('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "article"
	}
	return out
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func summarize(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) > 28 {
		words = words[:28]
	}
	return strings.Join(words, " ") + "..."
}

func inferPublishedDate(doc *goquery.Document, content string) (time.Time, string) {
	metaKeys := []string{
		`meta[property="article:published_time"]`,
		`meta[name="article:published_time"]`,
		`meta[name="pubdate"]`,
		`meta[name="publishdate"]`,
		`meta[name="date"]`,
	}
	for _, key := range metaKeys {
		if val, ok := doc.Find(key).Attr("content"); ok {
			if parsed, ok := parseDateCandidate(val); ok {
				return parsed, "meta"
			}
		}
	}

	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if parsed, ok := parseDateCandidate(line); ok {
			return parsed, "page_text"
		}
		if len(line) > 80 {
			continue
		}
		if match := monthDayYearPattern.FindString(line); match != "" {
			if parsed, ok := parseDateCandidate(match); ok {
				return parsed, "page_text"
			}
		}
		if match := monthYearPattern.FindString(line); match != "" {
			if parsed, ok := parseDateCandidate(match); ok {
				return parsed, "page_text"
			}
		}
	}

	return time.Time{}, "unknown"
}

func parseDateCandidate(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"January 2006",
		"Jan 2006",
		"1/2/2006",
		"1/2/06",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, raw)
		if err == nil && t.Year() >= 1990 && t.Year() <= time.Now().UTC().Year()+1 {
			if layout == "January 2006" || layout == "Jan 2006" {
				t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
			}
			return t.UTC(), true
		}
	}

	if match := monthDayYearPattern.FindString(raw); match != "" {
		if t, err := time.Parse("January 2, 2006", match); err == nil {
			return t.UTC(), true
		}
	}
	if match := monthYearPattern.FindString(raw); match != "" {
		if t, err := time.Parse("January 2006", match); err == nil {
			return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), true
		}
	}
	if match := shortDatePattern.FindString(raw); match != "" {
		for _, layout := range []string{"1/2/2006", "1/2/06"} {
			if t, err := time.Parse(layout, match); err == nil {
				return t.UTC(), true
			}
		}
	}

	return time.Time{}, false
}

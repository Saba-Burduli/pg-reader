package services

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestSlugify(t *testing.T) {
	got := slugify("The Lesson to Unlearn")
	if got != "the-lesson-to-unlearn" {
		t.Fatalf("unexpected slug: %s", got)
	}
}

func TestNormalizeText(t *testing.T) {
	got := normalizeText("hello \n\n world\t\t!")
	if got != "hello world !" {
		t.Fatalf("unexpected normalized text: %q", got)
	}
}

func TestInferPublishedDate(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><head><meta name="date" content="2019-03-14"/></head><body>x</body></html>`))
	if err != nil {
		t.Fatalf("failed to build test document: %v", err)
	}

	d, source := inferPublishedDate(doc, "Written on March 14, 2019 for founders.")
	if source != "meta" {
		t.Fatalf("expected meta source, got %s", source)
	}
	if d.Year() != 2019 || d.Month() != 3 || d.Day() != 14 {
		t.Fatalf("unexpected date: %v", d)
	}
}

func TestInferPublishedDateUnknown(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body>No date visible here.</body></html>`))
	if err != nil {
		t.Fatalf("failed to build test document: %v", err)
	}
	d, source := inferPublishedDate(doc, "No date visible here.")
	if !d.IsZero() {
		t.Fatalf("expected zero date, got: %v", d)
	}
	if source != "unknown" {
		t.Fatalf("expected unknown source, got %s", source)
	}
}

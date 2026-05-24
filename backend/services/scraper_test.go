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

	d := inferPublishedDate(doc, "", "Written on March 14, 2019 for founders.", "")
	if d.source != "meta" {
		t.Fatalf("expected meta source, got %s", d.source)
	}
	if d.sortTime.Year() != 2019 || d.sortTime.Month() != 3 || d.sortTime.Day() != 14 {
		t.Fatalf("unexpected date: %v", d)
	}
}

func TestInferPublishedDateUnknown(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body>No date visible here.</body></html>`))
	if err != nil {
		t.Fatalf("failed to build test document: %v", err)
	}
	d := inferPublishedDate(doc, "<html><body>No date visible here.</body></html>", "No date visible here.", "")
	if !d.sortTime.IsZero() {
		t.Fatalf("expected zero date, got: %v", d)
	}
	if d.source != "unknown" {
		t.Fatalf("expected unknown source, got %s", d.source)
	}
}

func TestInferPublishedDateFromTopArticleHTML(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><font size="2" face="verdana">May 2026<br><br>Essay text.</font></body></html>`))
	if err != nil {
		t.Fatalf("failed to build test document: %v", err)
	}
	d := inferPublishedDate(doc, `<html><body><img alt="Example" /><br><br><font size="2" face="verdana">May 2026<br><br>Essay text.</font></body></html>`, "Essay text.", "Example")
	if d.source != "page_text" {
		t.Fatalf("expected page_text source, got %s", d.source)
	}
	if d.sortTime.Year() != 2026 || d.sortTime.Month() != 5 || d.sortTime.Day() != 1 {
		t.Fatalf("unexpected date: %v", d)
	}
	if d.display != "2026-05" {
		t.Fatalf("unexpected display date: %s", d.display)
	}
}

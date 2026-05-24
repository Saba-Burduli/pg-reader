package services

import "testing"

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
	d, source := inferPublishedDate("Written on March 14, 2019 for founders.", 7)
	if source != "extracted" {
		t.Fatalf("expected extracted source, got %s", source)
	}
	if d.Year() != 2019 || d.Month() != 3 || d.Day() != 14 {
		t.Fatalf("unexpected date: %v", d)
	}
}

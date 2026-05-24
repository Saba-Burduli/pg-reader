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

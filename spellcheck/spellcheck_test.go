package spellcheck

import (
	"testing"
)

func TestSpellcheckPositive(t *testing.T) {
	input := "# This is a heading"
	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 0 {
		t.Fatal("Errs should be empty")
	}
}

func TestSpellcheckNeg(t *testing.T) {
	input := "# Thes is a heading"
	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 1 {
		t.Fatal("Errs should have one element")
	}
}

func TestSpellcheckDash(t *testing.T) {
	input := "# This-is-a-heading"
	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 0 {
		t.Fatal("Errs should have one element")
	}
}

func TestSpellcheckMultierr(t *testing.T) {
	input := "# Ths isy a headin"
	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 3 {
		t.Fatal("Errs incorect count")
	}
}

func TestSpellcheckMultiline(t *testing.T) {
	input := "# This is a heading\n> and this is a quote"
	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 0 {
		t.Fatal("Errs incorect count")
	}
}

func TestMultilineFileNoErrs(t *testing.T) {
	input := `# This is a heading

Then some body here **bold**, some _italic_

> A quote from a **smart** person

- This
- And that
- And other`

	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 0 {
		t.Fatal("Errs incorect count")
	}
}

func TestMultilineFileWithErrs(t *testing.T) {
	input := `# Ths isy a headin

Than somed body here **bolf**, some _italic_

> A quote from a **smort** person

- This
- And thay
- And other`

	s := NewChecker()
	errs := s.CheckSpelling(input)
	if len(errs) != 7 {
		t.Fatal("Errs incorect count")
	}
}

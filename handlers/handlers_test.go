package handlers

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"github.com/strjkc/noteapi/spellcheck"
)

var checker spellcheck.SpellChecker

func TestMain(m *testing.M) {
	fact := spellcheck.NewWordMapFactory(os.Getenv("DICTDIR"))
	wm, err := fact.WordMap("eng")
	if err != nil {
		return
	}
	checker = spellcheck.NewChecker(wm)
	m.Run()
}

func TestSpellCheck(t *testing.T) {
	parser := spellcheck.NewParser(checker)
	errors, err := parser.SpellCheckerService(multipart.NewReader(strings.NewReader(`# This is a heading

Then some body here **bold**, some _italic_

> A quote from a **smart** person

- This
- And that
- And other`), "\r\n"))
	if err != nil {
		t.Fatal("Error")
	}
	for _, error := range errors {
		if len(error.Mistakes) != 0 {
			t.Fatal("Errors not 0")
		}
	}
}

func TestSpellCheck2(t *testing.T) {
	parser := spellcheck.NewParser(checker)
	errors, err := parser.SpellCheckerService(multipart.NewReader(strings.NewReader(`# This is a heaing

Then some body here **bold**, some _italic_

> A quote from a **smart** person

- This
- And that
- And other`), "\r\n"))
	if err != nil {
		t.Fatal("Error")
	}
	for _, error := range errors {
		if error.Line == 1 {
			if len(error.Mistakes) != 1 {
				t.Fatal("Errors not 1")
			}
		} else {
			if len(error.Mistakes) != 0 {
				t.Fatal("Errors not 0")
			}
		}
	}
}

func TestSpellCheck3(t *testing.T) {
	parser := spellcheck.NewParser(checker)
	errors, err := parser.SpellCheckerService(multipart.NewReader(strings.NewReader(`# This is a heaing

Then some boxdy here **boldx**, some _italicqsee_

> A quote from a **smort** person

- This
- And that
- And other`), "\r\n"))
	if err != nil {
		t.Fatal("Error")
	}
	for _, error := range errors {
		if error.Line == 1 {
			if len(error.Mistakes) != 1 {
				t.Fatal("Errors not 1")
			}
		} else if error.Line == 3 {
			if len(error.Mistakes) != 3 {
				t.Fatal(fmt.Sprintf("Errors not 3 - %d", len(error.Mistakes)))
			}
		} else if error.Line == 5 {
			if len(error.Mistakes) != 1 {
				t.Fatal("Errors not 1")
			}
		} else {
			if len(error.Mistakes) != 0 {
				t.Fatal("Errors not 0")
			}
		}
	}
}

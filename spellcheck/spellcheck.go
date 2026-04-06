package spellcheck

import (
	"strings"
)

type SP struct {
	wordMap *WordMap
	input   string
	ch      byte
	pos     int
	readPos int
}

func NewChecker(wordMap *WordMap) *SP {
	sp := SP{}
	sp.wordMap = wordMap
	return &sp
}

func (s *SP) initState(input string) {
	s.input = input
	s.readPos = 0
	s.pos = s.readPos
	if len(s.input) == 0 {
		s.ch = byte(0)
	} else {
		s.ch = s.input[s.readPos]
	}
}

func (s *SP) CheckSpelling(input string) []string {
	s.initState(input)
	var errors []string
	s.read()
	for s.readPos < len(input) {
		if s.isChar() {
			word := s.readWord()
			if _, ok := s.wordMap.ValidWords[word]; !ok {
				errors = append(errors, word)
			}
		} else {
			s.read()
		}
	}

	return errors
}

func (s *SP) read() {
	if s.readPos < len(s.input) {
		s.ch = s.input[s.readPos]
	}
	s.pos = s.readPos
	s.readPos++
}

func (s *SP) isChar() bool {
	return s.ch >= 65 && s.ch <= 90 || s.ch >= 97 && s.ch <= 122
}

func (s *SP) readWord() string {
	pos := s.pos
	for s.isChar() && s.readPos <= len(s.input) {
		s.read()
	}
	return strings.ToLower(s.input[pos:s.pos])
}
